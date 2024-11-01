echo -e "{
  \"runtimeOptions\": {
    \"tfm\": \"net8.0\",
    \"framework\": {
      \"name\": \"Microsoft.NETCore.App\",
      \"version\": \"$(dotnet --info | awk '/Host:/ {getline; print $2}')\"
    }
  }
}" > /tmp/run.runtimeconfig.json

echo -e "
global using global::System;
global using global::System.Collections.Generic;
global using global::System.IO;
global using global::System.Linq;
global using global::System.Net.Http;
global using global::System.Threading;
global using global::System.Threading.Tasks;
" > /tmp/GlobalUsings.cs
