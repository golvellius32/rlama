#!/bin/bash
# Script de compilation pour RLAMA

VERSION=$(grep "Version = " cmd/root.go | cut -d'"' -f2)
PLATFORMS=("windows/amd64" "windows/386" "darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")
BINARY_NAME="rlama"

echo "Compilation de RLAMA v${VERSION}..."

rm -rf ./dist
mkdir -p ./dist

for platform in "${PLATFORMS[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$BINARY_NAME'_'$GOOS'_'$GOARCH
    
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Compilation pour $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o ./dist/$output_name
    
    if [ $? -ne 0 ]; then
        echo "Erreur lors de la compilation pour $GOOS/$GOARCH"
    else
        echo "Compilation réussie pour $GOOS/$GOARCH"
    fi
done

echo "Création des archives..."
cd ./dist
for file in rlama_*
do
    zip "${file}.zip" "$file"
done

echo "Compilation terminée!" 