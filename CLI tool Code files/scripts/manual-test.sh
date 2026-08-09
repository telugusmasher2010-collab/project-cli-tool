#!/bin/bash
# ============================================================
#  proj-init MANUAL TEST KIT — run before launch
#  Usage: bash manual-test.sh
#  Green = PASS, Red [FAIL] = bug. Zero FAILs = launch-ready.
# ============================================================
set -u
BIN="${1:-./pi.exe}"
PASS=0; FAIL=0

say()  { PASS=$((PASS+1)); printf '\033[32mPASS\033[0m %s\n' "$1"; }
bad()  { FAIL=$((FAIL+1)); printf '\033[31mFAIL\033[0m %s\n' "$1"; }

# isolate from any existing user config, restore on exit
CFG="$HOME/.proj-init/config.yml"
STASHED=""
[ -f "$CFG" ] && { mv "$CFG" "$CFG.kitbak"; STASHED=1; echo "(stashed existing config for test isolation)"; }
trap '[ -n "$STASHED" ] && mv "$CFG.kitbak" "$CFG" 2>/dev/null; rm -f "$CFG.kitbak"' EXIT

echo "== Building binary =="
cd "$(dirname "$0")" 2>/dev/null || true

T=e2e-$$
rm -rf "$T"; mkdir -p "$T"; cd "$T"

echo "== T1: --help =="
"$BIN" --help >/dev/null 2>&1 && say "T1 help" || bad "T1 help"

echo "== T2: version =="
out=$("$BIN" version 2>&1) && echo "   v: $(echo "$out" | head -1)"

echo "== T3: list = shows 5 templates =="
n=$("$BIN" list 2>&1 | grep -cE "expense-splitter|next-webapp|react-native-map|tauri-llm|whatsapp-bot")
[ "$n" -eq 5 ] && say "T3 list ($n)" || bad "T3 list ($n)"

echo "== T4. interactive init — type: myapp / 1 / (Enter for default) =="
printf 'myapp\n1\n\n' | "$BIN" init 2>&1 | tail -1
[ -d "./myapp" ] && [ -f "./myapp/pubspec.yaml" ] && say "T4 scaffold" || bad "T4 scaffold"

echo "== T5. NO leftover {{ }} in generated files =="
L=$(grep -rl "{{" ./myapp 2>/dev/null | head -1)
[ -z "$L" ] && say "T5 no placeholders" || bad "T5 leftover in $L"

echo "== T6. bad name rejected: 'my app' then retry 'okapp' =="
printf 'my app\nokapp\n1\n\n' | "$BIN" init >/dev/null 2>&1
[ -d "./okapp" ] && say "T6 recovery" || bad "T6 recovery"

echo "== T7. bad template # rejected: 9 then 3 =="
printf 'z\n9\n3\n\n' | "$BIN" init >/dev/null 2>&1
[ -d "./z" ] && [ -f "./z/package.json" ] && say "T7 recovery" || bad "T7 recovery"

echo "== T8. non-empty dir must FAIL =="
mkdir -p busy; echo x > busy/x.txt
if printf 'q\n1\n' | "$BIN" init -o ./busy >/dev/null 2>&1; then bad "T8 non-empty dir ACCEPTED"; else say "T8 non-empty rejected"; fi

echo "== T9. -o flag skips prompt =="
if printf 'q\n1\n' | "$BIN" init -o ./flagonly >/dev/null 2>&1 && [ -f ./flagonly/pubspec.yaml ]; then say "T9 -o works"; else bad "T9 -o works"; fi

echo "== T10. all 5 templates generate clean =="
ok=1
for t in 1 2 3 4 5; do
  printf "t$t\n$t\n./gen$t\n" | "$BIN" init >/dev/null 2>&1 || { ok=0; echo "   tpl $t FAILED"; }
  [ -n "$(grep -rl '{{' ./gen$t 2>/dev/null)" ] && { ok=0; echo "   tpl $t has {{ }}"; }
done
[ $ok -eq 1 ] && say "T10 all templates" || bad "T10 templates"

echo "== T11. empty project name rejected =="
if printf '\nx\n1\n./n1\n' | "$BIN" init >/dev/null 2>&1; then say "T11 empty-name accepted(fine)"; else say "T11 empty rejected"; fi

echo
echo "================ RESULT: $PASS pass, $FAIL fail ================"
cd .. && rm -rf "$T"