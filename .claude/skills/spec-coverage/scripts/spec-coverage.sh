#!/usr/bin/env bash
# Fetch the published Waroom OpenAPI spec and compare its operations with the ones wrapped in waroom/.
#
# Usage: spec-coverage.sh [<output path for the fetched spec>]
#   OPENAPI_URL overrides the spec URL.
#
# Output (TSV, one line per operation in the published spec or in api/openapi.yaml):
#   status  tag  method  path  operationId  wrapper  summary
# status:
#   wrapped    in the published spec and called from waroom/ (wrapper = Service.Method)
#   missing    in the published spec, not called from waroom/
#   spec-only  in api/openapi.yaml (and the published spec) but not called from waroom/ — generated but not wrapped
#   removed    in api/openapi.yaml but no longer in the published spec
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

url=${OPENAPI_URL:-https://static.waroom.com/openapi/v0.yaml}
out=${1:-$(mktemp -t waroom-openapi.XXXXXX).yaml}

curl -sfSL --max-time 30 "$url" -o "$out"

# operations: operationId<TAB>tag<TAB>method<TAB>path<TAB>summary
# Relies on the layout of the published spec: paths at 2 spaces, methods at 4, fields at 6.
list_ops() {
  awk '
    function flush() { if (op != "") print op "\t" tag "\t" method "\t" path "\t" summary; op = "" }
    /^paths:/ { inpaths = 1; next }
    /^[^ ]/   { flush(); inpaths = 0 }
    !inpaths  { next }
    /^  [^ ]/ { flush(); path = $1; sub(/:$/, "", path); gsub(/"/, "", path); next }
    /^    [^ ]/ { flush(); method = $1; sub(/:$/, "", method); method = toupper(method); tag = ""; summary = ""; next }
    /^      tags:/ { want_tag = 1; next }
    want_tag && /^      - / { if (tag == "") tag = substr($0, 9); next }
    { want_tag = 0 }
    /^      summary:/ { summary = substr($0, 16); gsub(/^"|"$/, "", summary) }
    /^      operationId:/ { op = $2 }
    END { flush() }
  ' "$1"
}

# wrapped: OperationId (PascalCase)<TAB>Service.Method
list_wrapped() {
  awk '
    /^func \(s \*[A-Za-z]+Service\) [A-Z]/ {
      recv = $3; sub(/^\*/, "", recv); sub(/\)$/, "", recv)
      name = $4; sub(/\(.*/, "", name)
      fn = recv "." name
    }
    /s\.client\.raw\.[A-Z][A-Za-z]+\(/ {
      match($0, /raw\.[A-Z][A-Za-z]+\(/)
      print substr($0, RSTART + 4, RLENGTH - 5) "\t" fn
    }
  ' $(ls waroom/*.go | grep -v '_test.go$') | sort -u
}

spec_version() { awk '/^info:/{f=1;next} f&&/^[^ ]/{exit} f&&/^  version:/{v=$2; gsub(/["\047]/,"",v); print v; exit}'; }

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

list_ops "$out" > "$tmp/pub"
list_ops api/openapi.yaml > "$tmp/local"
list_wrapped > "$tmp/wrapped"

echo "# spec: $url"
echo "# fetched: $out"
echo "# version: published $(spec_version < "$out") / repo spec $(spec_version < api/openapi.yaml)"
echo "# repo spec: api/openapi.yaml"
echo "# operations: published $(wc -l < "$tmp/pub" | tr -d " ") / repo spec $(wc -l < "$tmp/local" | tr -d " ") / wrapped $(cut -f1 "$tmp/wrapped" | sort -u | wc -l | tr -d " ")"
echo
printf 'status\ttag\tmethod\tpath\toperationId\twrapper\tsummary\n'

awk -F'\t' -v OFS='\t' '
  function pascal(s) { return toupper(substr(s, 1, 1)) substr(s, 2) }
  FNR == 1 { file++ }
  file == 1 { w[$1] = (w[$1] == "" ? $2 : w[$1] "," $2); next }
  file == 2 { local[$1] = 1; next }
  file == 3 {
    pub[$1] = 1
    p = pascal($1)
    if (p in w)          st = "wrapped"
    else if ($1 in local) st = "spec-only"
    else                 st = "missing"
    print st, $2, $3, $4, $1, (p in w ? w[p] : "-"), $5
    next
  }
  file == 4 && !($1 in pub) {
    p = pascal($1)
    print "removed", $2, $3, $4, $1, (p in w ? w[p] : "-"), $5
  }
' "$tmp/wrapped" "$tmp/local" "$tmp/pub" "$tmp/local" | sort -t$'\t' -k1,1 -k2,2 -k4,4
