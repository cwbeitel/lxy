import argparse
import random

'''

Given a sorted vcf file, generate a version of it with a field
tagging certain variants as belonging to the same haplotype block
(being in phase). Variant blocks that are not in the same phase
should be randomly inverted relative to eachother.

python /Users/cb/code/src/github.com/cb01/core/lxy/scripts/simvarblocks.py \
--infile_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/vars/bychr/vars.1.vcf \
--outfile_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/vars/bychr/blocks/vars.1.blk.vcf

python /Users/cb/code/src/github.com/cb01/core/lxy/scripts/simvarblocks.py \
--infile_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/vars/bychr/vars.22.vcf \
--outfile_path=/Users/cb/code/src/github.com/cb01/core/lxy/data/GM12878/vars/bychr/blocks/vars.22.blk.vcf

'''

if __name__ == '__main__':

	parser = argparse.ArgumentParser(description='Split a gzipped sam file by a specified field')
	parser.add_argument('--infile_path',help='Path to input file')
	parser.add_argument('--outfile_path',help='Output directory')
	parser.add_argument('--bsize_min',help='Output directory', default=100000)
	parser.add_argument('--bsize_max',help='Output directory', default=100000)
	parser.add_argument('--alternate',help='Whether to alternate the phase of the variant blocks.', default=False)
	args = parser.parse_args()

	current_block_size = -1
	current_block_size_target = random.randint(args.bsize_min, args.bsize_max)
	invert_state = 0
	block_counter = 1

	#llim = 5000
	#lcount = 0

	out = open(args.outfile_path, 'w')

	with open(args.infile_path, 'r') as f:

		for line in f:

			if line[0] != "#":

				current_block_size += 1

				#lcount += 1
				#if lcount > llim:
				#	break

				# Get the fields of the line
				arr = line.split()

				# If we do not have a current block
				if current_block_size >= current_block_size_target:

					current_block_size = 0

					# Choose a new current block size target
					current_block_size_target = random.randint(args.bsize_min, args.bsize_max)

					# Choose a new invert state
					#invert_state = random.randint(0,1)
					if invert_state == 1:
						invert_state = 0
					else:
						invert_state = 1

					# Increment the block counter
					block_counter += 1

				# Update the block id field
				arr[7] = arr[7] + ";BLOCK=" + str(block_counter)

				if args.alternate:
					# If the invert state is 1, invert ref and alt
					if invert_state == 1:
						ref = arr[3]
						alt = arr[4]
						arr[3] = alt
						arr[4] = ref

				# Write the output line
				out.write('\t'.join(arr)+"\n")

	out.close()



