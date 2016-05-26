import argparse
import gzip
import os

'''

python /Users/cb/code/src/github.com/cb01/core/lxy/scripts/split.py \
--infile_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/hic/SRR927086.nomaskref.sam.gz \
--out_base_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/hic/split

python /Users/cb/code/src/github.com/cb01/core/lxy/scripts/split.py \
--infile_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/hic/test.sam.gz \
--out_base_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/hic/split

needs to be able to match a string to ignore in the field as well, for example not generate an output
file named *.sam

prefer to write output in non-gzipped form

'''

if __name__ == '__main__':

	parser = argparse.ArgumentParser(description='Split a gzipped sam file by a specified field')
	parser.add_argument('--infile_path',help='Path to input file')
	parser.add_argument('--out_base_path',help='Output directory')
	parser.add_argument('--ignore_char',help='',default='@')	
	parser.add_argument('--split_column',help='',default='3',type=int)
	args = parser.parse_args()

	files = {}

	with gzip.open(args.infile_path, 'rb') as f:

		for line in f:

			# If the line does not start with the ignore character
			if line[0] is not args.ignore_char:

				# Get the value of the field in question, v
				v = line.split()[args.split_column-1]

				if v is not "*":
					# If v is not a key in the files map, create it and open the file
					if v not in files:
						files[v] = open(os.path.join(args.out_base_path, v + ".sam"), 'w')

					# Get the file handle from the files map that corresponds to v
					fh = files[v]

					# Write the current line to the output file
					fh.write(line)

	# Close each file in the files map
	for fh in files.values():
		fh.close()
