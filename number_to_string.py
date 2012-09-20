DICT = { str(i+1):chr((ord('a')+ i)) for i in range(26) }
MAX_SIZE = max(len(k) for k in DICT.keys())

def parse(dt,cache,string):

	if string == "":
		return [""]

	if string in cache:
		return cache[string]

	ret = []
	for i in range(MAX_SIZE):
		if len(string) > i and string[0:i+1] in dt:
			ret += [ dt[string[0:i+1]] + j for j in parse(dt,cache,string[i+1:]) ]
		else:
			break

	cache[string] = ret
	return ret

if __name__ == '__main__':

	cache = {}
	ret = parse(DICT,cache,"")
	print ret

	ret = parse(DICT,cache,"12259")
	print ret

	# ret = parse(dt,cache,"12323543267536423543265362354374342345276352435165363246743753")
	# print ret


