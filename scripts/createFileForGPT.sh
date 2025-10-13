#!/bin/bash
# bash scripts/createFileForGPT.sh

OUTPUT_FILE=~/projectForGPT.txt
EXCLUDE_DIRS=(".git" "pg_data" "configs")
EXCLUDE_FILES=("go.sum" "go.mod" "*.log" "*.tmp" ".env")

# Очистка/создание файла вывода
echo "--- Дерево файлов ---" > "$OUTPUT_FILE"
tree >> "$OUTPUT_FILE"

# Формируем выражения для find
build_exclude_expr() {
  local type="$1"
  shift
  local expr=()
  for name in "$@"; do
    expr+=("-name" "$name" "-o")
  done
  (( ${#expr[@]} )) && unset 'expr[-1]' # удалить последний "-o"
  echo "(" -type "$type" "${expr[@]}" ")" -prune -o
}

# Строим команду find
FIND_CMD=("find" ".")

# Добавляем исключения для директорий
if [[ ${#EXCLUDE_DIRS[@]} -gt 0 ]]; then
  IFS=' ' read -r -a exclude_dirs_expr <<< "$(build_exclude_expr d "${EXCLUDE_DIRS[@]}")"
  FIND_CMD+=("${exclude_dirs_expr[@]}")
fi

# Ищем только файлы
FIND_CMD+=("-type" "f")

# Добавляем исключения для файлов
if [[ ${#EXCLUDE_FILES[@]} -gt 0 ]]; then
  # Создаём отдельное выражение для исключений файлов
  EXPR=()
  for file in "${EXCLUDE_FILES[@]}"; do
    EXPR+=("-name" "$file" "-o")
  done
  unset 'EXPR[-1]' # удаляем последний "-o"

  # Добавляем инвертированное условие
  FIND_CMD+=("!" "(" "${EXPR[@]}" ")")
fi

# Добавляем выполнение cat
FIND_CMD+=("-exec" "sh" "-c" 'printf "\n\n\n--- Содержимое файла %s ---\n" "$1"; cat "$1"' "_" "{}" ";")

#echo "Выполняется команда:"
#printf "%q " "${FIND_CMD[@]}"; echo

# Выполняем команду
"${FIND_CMD[@]}" >> "$OUTPUT_FILE"
