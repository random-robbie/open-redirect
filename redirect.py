import sys
import subprocess
import shlex
import re
import os
from time import sleep

CASPERJS_SCRIPT = "redirect.js"
XVFB_RUN = "/usr/bin/xvfb-run"

if len(sys.argv) < 3:
	sys.exit("Usage: python redirect.py urls.txt payloads.txt")

urllist = sys.argv[1]
payloadlist = sys.argv[2]



def log_results(URL, result):
	print("\n\n\n\n\n\n[*]*****Open Redirect Found*****[*]")
	print("[*]" + URL + "[*]")
	print("[*]Redirects to: " + result + "[*]")
	print("\n\n\n\n")
	with open("found.txt", "a") as file:
		file.write(URL + "\n") 

def test_redirect(url, payload):
	exmp = url.replace("https://", "")
	exmp = exmp.replace("http://", "")
	payload = payload.replace("example.com", "redirect.xsses.rocks")
	URL = url + payload

	if "DOCKER" in os.environ:
		cmd = f"casperjs {CASPERJS_SCRIPT}"
	else:
		cmd = f"{XVFB_RUN} -a casperjs {CASPERJS_SCRIPT}"

	args = shlex.split(cmd)
	args.append(URL)
	stdout_as_string = subprocess.check_output(args).decode('utf-8')
	result = re.sub(r'\x1b\[[\d;]+m', '', stdout_as_string)
	result = result.strip()
	found = False

	if result.startswith(("http://google.com", "https://google.com", "http://redirect.xsses.rocks", "https://redirect.xsses.rocks")):
		log_results(URL, result)
		found = True

	if not found:
		print(f"[*]Open Redirect NOT Found for {URL}[*]")


def main(urllist, payloadlist):
	with open(urllist) as f:
		os.system('clear')
		print("\n\n\n\n")
		print("[*] ***************************************[*]")
		print("[*] Open Redirect Finder By @Random_Robbie [*]")
		print("[*] ***************************************[*]")
		print("\n\n")
		sleep(2)
		print("[*] Searching for Open Redirects [*]\n\n")

		for line in f:
			line = line.strip()
			url = line

			with open(payloadlist) as g:
				for payload in g:
					payload = payload.strip()
					try:
						test_redirect(url, payload)
					except Exception as e:
						print(f"[*]Debug: {url} & Payload {payload} [*]")
						print(f'Error: {e}')


try:
	main(urllist, payloadlist)
	print("[*]Finished Every Payload.... No Open Redirect Found[*]")
except Exception as e:
	print(e)
	print(f'Error: {e}')
	sys.exit(1)
