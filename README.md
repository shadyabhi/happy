# Happy

A small program to analyze who wins the race between ipv4 and ipv6.

	➜ $?=0 @arastogi-mn3 golang/happy [ 9:41AM] (master|…)➤ go run main.go -server 'www.google.com:443' -timeout 1000
	2017/11/05 09:41:31 Connecting to server: www.google.com:443 with timeout: 1000
	2017/11/05 09:41:31 Connected to address: 216.58.203.228:443 in 63ms
	2017/11/05 09:41:31 Connected to address: [2404:6800:4003:802::2004]:443 in 65ms

	>>>  1s elasped...
	➜ $?=0 @arastogi-mn3 golang/happy [ 9:41AM] (master|…)➤ go run main.go -server 'www.facebook.com:443' -timeout 1000
	2017/11/05 09:41:36 Connecting to server: www.facebook.com:443 with timeout: 1000
	2017/11/05 09:41:36 Connected to address: 157.240.16.35:443 in 43ms
	2017/11/05 09:41:36 Connected to address: [2a03:2880:f12f:83:face:b00c:0:25de]:443 in 46ms

	>>>  0s elasped...
	➜ $?=0 @arastogi-mn3 golang/happy [ 9:41AM] (master|…)➤
