#!/bin/bash
# bash scripts/createFileForGPT.sh

OUTPUT_FILE=projectForGPT.txt
EXCLUDE_DIRS=(".git" "pg_data" "config" "bin" ".idea" "static" "docs/static/swagger")
EXCLUDE_FILES=(OUTPUT_FILE "go.sum" "go.mod" "*.log" "*.tmp" ".env" "000002_cities.up.sql")

# Очистка/создание файла вывода
echo "--- Дерево файлов ---" > "$OUTPUT_FILE"
tree >> "$OUTPUT_FILE"

touch "$OUTPUT_FILE"
> "$OUTPUT_FILE"

echo "--- Дерево файлов ---" > "$OUTPUT_FILE"
tree >> "$OUTPUT_FILE"

# Формируем команду find
FIND_CMD=(find . -type f)


for dir in "${EXCLUDE_DIRS[@]}"; do
	FIND_CMD+=(-not -path "*/${dir}/*" -not -path "./${dir}")
done

for pattern in "${EXCLUDE_FILES[@]}"; do
	FIND_CMD+=(-not -name "$pattern")
done


FIND_CMD+=(-exec sh -c '
file="$1"
printf "\n\n\n--- Содержимое файла %s ---\n" "$file"
cat "$file"
' _ {} \;)

"${FIND_CMD[@]}" >> "$OUTPUT_FILE"

echo "cat $OUTPUT_FILE"