import sys
import subprocess
import shlex
import re
import os
from time import sleep

CASPERJS_SCRIPT = "redirect.js" 
XVFB_RUN = "/usr/bin/xvfb-run" #Location of XVFB to run headless

if len(sys.argv) < 2:
		sys.exit("Usage: python redirect.py uber.txt payloads.txt")
		exit();
		
		
urllist = sys.argv[1]
payloadlist = sys.argv[2]



def log_results(URL,result):
		print "\n\n\n\n\n\n[*]*****Open Redirect Found*****[*]"
		print "[*]"+URL+"[*]"
		print "[*]Redirects to: "+result+"[*]"
		print "\n\n\n\n"
		file = open("found.txt","a") 
		file.write(""+URL+"\n")
		file.close() 

def test_redirect(url,payload):

	exmp = url.replace("https://","")
	exmp = exmp.replace("http://","")
	payload = payload.replace("example.com","redirect.xsses.rocks")
	URL = ""+url+""+payload+""
	cmd = ""+XVFB_RUN+" -a casperjs "+CASPERJS_SCRIPT+""
	args = shlex.split(cmd)
	args.append(URL)
	stdout_as_string = subprocess.check_output(args).decode('utf-8')
	result = re.sub(r'\x1b\[[\d;]+m', '', stdout_as_string)
	result = result.encode('ascii','ignore')
	result = result.replace("\n","")
	found = False
	if result.startswith(("http://google.com","https://google.com","http://redirect.xsses.rocks","https://redirect.xsses.rocks")) == True:
		log_results(URL,result)
		found = True
	
	

	if found == False:
		print "[*]Open Redirect NOT  Found for "+URL+"[*]"

		

		
def main(urllist,payloadlist):

	with open(urllist) as f:
		os.system('clear')
		print "\n\n\n\n"
		print "[*] ***************************************[*]"
		print "[*] Open Redirect Finder By @Random_Robbie [*]"
		print "[*] ***************************************[*]"
		print "\n\n"
		sleep (2)
		print "[*] Searching for Open Redirects [*]\n\n"
		for line in f:
			line = line.replace("\r\n","")
			line = line.replace ("\n","")
			url = line
			with open(payloadlist) as g:
				for payload in g:
					payload = payload.replace("\r\n","")
					payload = payload.replace ("\n","")
					try:
						test_redirect(url,payload)
					except Exception as e:
						pass
						print "[*]Debug: "+url+" & Payload "+payload+" [*]"
						print('Error: %s' % e)
						pass
			
						
						
try:
	main(urllist,payloadlist)
	print "[*]Finished Every Payload.... No Open Redirect Found[*]"
except Exception as e:
		print (e)
		print('Error: %s' % e)
		sys.exit(1)
