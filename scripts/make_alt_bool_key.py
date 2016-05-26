import sys

state = 1
for i in range(1, int(sys.argv[1])+1):
	if state == 1:
		state = 0
		print sys.argv[3] + "_" + str(i), 1
	else:
		state = 1
		print sys.argv[3] + "_" + str(i), 0
