import requests
import sys
from time import sleep
requests.packages.urllib3.disable_warnings()

session = requests.Session()


if len(sys.argv) < 2:
		sys.exit("Usage: python redirect.py uber.txt payloads.txt")
		exit();
		
		
urllist = sys.argv[1]
payloadlist = sys.argv[2]

def test_redirect(url,payload):

	headers = {"Accept":"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8","User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0","Connection":"close","Accept-Language":"en-US,en;q=0.5","Accept-Encoding":"gzip, deflate"}
	exmp = url.replace("https://","")
	exmp = exmp.replace("http://","")
	payload = payload.replace("example.com",""+exmp+"")
	URL = ""+url+""+payload+""
	response = session.get(URL, headers=headers, verify=False)
	try:
		if response.status_code == 303:
			print "\n\n\n\n\n\n[*]*****Open Redirect Found*****[*]"
			print "[*]"+URL+"[*]"
			print "\n\n\n\n"
			file = open("found.txt","a") 
			file.write(""+URL+"\n")
			file.close() 
		
		if response.status_code == 302:
			if payload in response.history[0].headers['Location']:
				print "\n\n\n\n\n\n[*]*****Open Redirect Found*****[*]"
				print "[*]"+URL+"[*]"
				print "\n\n\n\n"
				file = open("found.txt","a") 
				file.write(""+URL+"\n")
				file.close() 


		else:
			print "[*]Open Redirect NOT  Found for "+URL+" HTTP ["+str(response.status_code)+"][*]"
	except requests.exceptions.MissingSchema:
		print("*] You forgot the protocol. http://, https://, [*]")
	except requests.exceptions.ConnectionError:
		print("*] Sorry, but I couldn't connect. There was a connection problem. [*]")
		pass
	except requests.exceptions.Timeout:
		print("*] Sorry, but I couldn't connect. I timed out. [*]")
		pass
	except requests.exceptions.TooManyRedirects:
		print("*] There were too many redirects.  I can't count that high. [*]")
		pass
	except Exception as e:
		print("*] Debug: "+e+" [*]")
		pass

		

		
def main(urllist,payloadlist):

	with open(urllist) as f:
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
