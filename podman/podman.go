//
//  Copyright 2024 whipcode.app (AnnikaV9)
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on
//  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific
//  language governing permissions and limitations under the License.
//

package podman

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/karlseguin/ccache/v3"
)

func NewExecutor(timeout int) *Executor {
	cache := ccache.New(ccache.Configure[map[string]interface{}]().MaxSize(100).ItemsToPrune(10))
	return &Executor{execCache: cache, timeout: timeout}
}

func (ex *Executor) RunCode(code, entry, ext, img string, enableCache bool) (int, map[string]interface{}) {
	if enableCache {
		if item := ex.execCache.Get(code); item != nil {
			go item.Extend(time.Hour * 24)
			return http.StatusOK, item.Value()
		}
	}

	boxID := strconv.Itoa(rand.Intn(9000000) + 1000000)
	srcFileName := fmt.Sprintf("run%s.%s", boxID, ext)
	srcFilePath := filepath.Join("./run", srcFileName)

	if err := os.WriteFile(srcFilePath, []byte(code), 0644); err != nil {
		log.Println("Could not write temp file:", err)
		return http.StatusInternalServerError, map[string]interface{}{
			"detail": "internal server error",
		}
	}
	defer os.Remove(srcFilePath)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ex.timeout)*time.Second)
	defer cancel()

	var stdout, stderr bytes.Buffer
	args := []string{
		"run",
		"--rm",
		"--read-only",
		"--no-hosts",
		"--hostname", fmt.Sprintf("box%s", boxID),
		"--network", "none",
		"--timeout", fmt.Sprintf("%d", ex.timeout+2),
		"--cap-drop", "ALL",
		"--memory", "512m",
		"--memory-reservation", "128m",
		"--cpus", "1.0",
		"--pids-limit", "32",
		"--user", "nobody",
		"--tmpfs", "/tmp:rw,size=64m,mode=1777",
		"--tmpfs", "/var/tmp:ro,size=32m,mode=1777",
		"--security-opt", "no-new-privileges",
		"--security-opt", "mask=/home:/etc:/opt:/media:/root:/run:/srv:/sys:/var",
		"--security-opt", "label=type:whipcode.process",
		"--volume", fmt.Sprintf("./entry/%s.sh:/entry.sh:z,ro", entry),
		"--volume", fmt.Sprintf("./run/%s:/source.%s:Z,ro", srcFileName, ext),
		img, "sh", "-c", "echo stdout-start && echo stderr-start >&2 && sh ./entry.sh",
	}
	cmdExec := exec.CommandContext(ctx, "/usr/bin/podman", args...)
	cmdExec.Stdout = &stdout
	cmdExec.Stderr = &stderr

	startTime := time.Now()
	_ = cmdExec.Run()
	duration := time.Since(startTime).Seconds()

	if ctx.Err() == context.DeadlineExceeded {
		result := map[string]interface{}{
			"stdout":        "",
			"stderr":        "",
			"container_age": duration,
			"timeout":       true,
		}

		if enableCache {
			go ex.execCache.Set(code, result, time.Hour*24)
		}

		return http.StatusOK, result
	}

	stdoutStr := stdout.String()
	stderrStr := stderr.String()
	if !(strings.HasPrefix(stdoutStr, "stdout-start")) || !(strings.HasPrefix(stderrStr, "stderr-start")) {
		log.Println("Blocked unsafe output:", "stdout:", stdoutStr, "stderr:", stderrStr)
		return http.StatusInternalServerError, map[string]interface{}{
			"detail": "internal server error",
		}
	}

	result := map[string]interface{}{
		"stdout":        strings.TrimPrefix(stdoutStr, "stdout-start\n"),
		"stderr":        strings.TrimPrefix(stderrStr, "stderr-start\n"),
		"container_age": duration,
		"timeout":       false,
	}

	if enableCache {
		go ex.execCache.Set(code, result, time.Hour*24)
	}

	return http.StatusOK, result
}
