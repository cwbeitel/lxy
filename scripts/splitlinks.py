import argparse
import os

'''
python /Users/cb/code/src/github.com/cb01/core/lxy/scripts/splitlinks.py \
--infile_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/links/GM.1mbp.links \
--out_base_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/links/split
'''

if __name__ == '__main__':

	parser = argparse.ArgumentParser(description='Split links file by a specified field')
	parser.add_argument('--infile_path',help='Path to input file')
	parser.add_argument('--out_base_path',help='Output directory')
	parser.add_argument('--ignore_char',help='Output directory', default='#')
	args = parser.parse_args()

	files = {}

	with open(args.infile_path, 'rb') as f:

		for line in f:

			# If the line does not start with the ignore character
			if line[0] is not args.ignore_char:

				# Get the value of the field in question, v
				c1, c2 = line.split()[:2]
				c1base = c1.split("_")[0]
				c2base = c2.split("_")[0]
				if c1base == c2base :

					# If v is not a key in the files map, create it and open the file
					if c1base not in files:
						files[c1base] = open(os.path.join(args.out_base_path, c1base + ".links"), 'wb')

					# Get the file handle from the files map that corresponds to v
					fh = files[c1base]

					# Write the current line to the output file
					fh.write(line)

	# Close each file in the files map
	for fh in files.values():
		fh.close()


