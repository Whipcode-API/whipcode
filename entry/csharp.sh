REF_PATH=$(echo /usr/lib/dotnet/shared/Microsoft.NETCore.App/*)

REFERENCES=""

for dll in "$REF_PATH"/*.dll; do
    if [ -e "$dll" ]; then
        REFERENCES="$REFERENCES -r:$dll"
    fi
done

dotnet /usr/lib/dotnet/sdk/*/Roslyn/bincore/csc.dll //source.cs /tmp/GlobalUsings.cs -nologo -unsafe -out:/tmp/run.dll $REFERENCES && dotnet /tmp/run.dll
