#!/usr/bin/env python3
"""Budget plan smoke tester.

The script pushes CBOR payloads to the advisor endpoints. Two modes exist:

1. Default payload mode (no flags) sends the hard-coded sample data to
   ``/advisor/budget-plan/test``.
2. Username mode (``--username``) sends a lightweight body that only includes
   the username plus optional goal knobs, which is handy for the new
   history-backed route (e.g. ``/advisor/budget-plan/history``).

Usage:
    python scripts/test_budget_plan.py \
        --base-url http://localhost:8080 \
        --token "$OPENAI_API_KEY" \
        [--username jake --path /advisor/budget-plan/history]
"""

from __future__ import annotations

import argparse
import json
import os
import sys
from typing import Any, Dict

import cbor2
import requests


def build_payload() -> Dict[str, Any]:
  """Provide a sane default payload; tweak as needed."""
  return {
      "currency": "KRW",
      "expenses": [
          {"category": "식비", "amount": 420_000, "notes": "배달"},
          {"category": "교통", "amount": 80_000},
          {"category": "구독", "amount": 35_000, "notes": "OTT + 음악"},
      ],
      "target_goal": "여름 휴가비 50만원 마련",
      "desired_reduction_percent": 20,
  }


def build_username_payload(args: argparse.Namespace) -> Dict[str, Any]:
  """Build the minimal payload expected by the history-backed endpoint."""
  payload: Dict[str, Any] = {
      "username": args.username,
  }

  if args.currency:
    payload["currency"] = args.currency

  if args.target_goal:
    payload["target_goal"] = args.target_goal

  if args.desired_reduction_percent is not None:
    payload["desired_reduction_percent"] = args.desired_reduction_percent

  return payload


def main() -> None:
  parser = argparse.ArgumentParser(description="Hit /advisor/budget-plan/test with CBOR payloads")
  parser.add_argument("--base-url", default="http://localhost:8080", help="API origin (default: %(default)s)")
  parser.add_argument("--path", default="/advisor/budget-plan/test", help="Relative endpoint path")
  parser.add_argument("--payload-json", help="Optional JSON file overriding the default payload")
  parser.add_argument("--token", help="OpenAI API key; forwarded as header if provided")
  parser.add_argument("--username", help="Username to look up via the history-backed endpoint")
  parser.add_argument("--currency", default="KRW", help="Currency to use when --username is set")
  parser.add_argument("--target-goal", help="Optional goal text when --username is set")
  parser.add_argument(
      "--desired-reduction-percent",
      type=float,
      help="Optional reduction percentage when --username is set",
  )
  parser.add_argument(
      "--timeout",
      type=float,
      default=30,
      help="How many seconds to wait for the API response (default: %(default)s)",
  )

  args = parser.parse_args()

  if args.username and args.payload_json:
    parser.error("--username cannot be combined with --payload-json")

  if args.username:
    payload = build_username_payload(args)
  elif args.payload_json:
    with open(args.payload_json, "r", encoding="utf-8") as fp:
      payload = json.load(fp)
  else:
    payload = build_payload()

  url = f"{args.base_url.rstrip('/')}{args.path}"
  headers = {
      "Content-Type": "application/cbor",
      "Accept": "application/cbor",
  }

  # Helpful when you run this through a proxy or want to propagate the key
  token = args.token or os.getenv("OPENAI_API_KEY")
  if token:
    headers["X-OpenAI-Key"] = token

  try:
    response = requests.post(url, data=cbor2.dumps(payload), headers=headers, timeout=args.timeout)
  except requests.exceptions.Timeout:
    print(f"Request to {url} timed out after {args.timeout} seconds.")
    print("Try increasing --timeout or confirm the API is healthy.")
    sys.exit(1)
  except requests.exceptions.RequestException as exc:
    print(f"Request to {url} failed: {exc}")
    sys.exit(1)

  print(f"Status: {response.status_code}")

  if response.status_code >= 400:
    print(response.text)
    sys.exit(1)

  try:
    decoded = cbor2.loads(response.content)
    print(json.dumps(decoded, ensure_ascii=False, indent=2))
  except Exception as exc:  # pylint: disable=broad-except
    print(f"Failed to decode CBOR: {exc}")
    print(response.content)
    sys.exit(1)


if __name__ == "__main__":
  main()
