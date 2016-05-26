#!/usr/bin/env bash

# download reference (todo)

wget ftp://ftp.1000genomes.ebi.ac.uk/vol1/ftp/technical/reference/phase2_reference_assembly_sequence/hs37d5.fa.gz


# download variants (todo)

for c in 1 5 9 13 17; do wget -O - ftp://ftp-trace.ncbi.nih.gov/1000genomes/ftp/release/20130502/ALL.chr$c.phase3_shapeit2_mvncall_integrated_v5a.20130502.genotypes.vcf.gz | zcat - | cut -f1,2,3,4,5,6,7,8,1762 | gzip -9 > vars.$c.vcf.gz; done
for c in 2 6 10 14 18; do wget -O - ftp://ftp-trace.ncbi.nih.gov/1000genomes/ftp/release/20130502/ALL.chr$c.phase3_shapeit2_mvncall_integrated_v5a.20130502.genotypes.vcf.gz | zcat - | cut -f1,2,3,4,5,6,7,8,1762 | gzip -9 > vars.$c.vcf.gz ; done
for c in 3 7 11 15 19 21; do wget -O - ftp://ftp-trace.ncbi.nih.gov/1000genomes/ftp/release/20130502/ALL.chr$c.phase3_shapeit2_mvncall_integrated_v5a.20130502.genotypes.vcf.gz | zcat - | cut -f1,2,3,4,5,6,7,8,1762 | gzip -9 > vars.$c.vcf.gz ; done
for c in 4 8 12 16 20 22; do wget -O - ftp://ftp-trace.ncbi.nih.gov/1000genomes/ftp/release/20130502/ALL.chr$c.phase3_shapeit2_mvncall_integrated_v5a.20130502.genotypes.vcf.gz | zcat - | cut -f1,2,3,4,5,6,7,8,1762 | gzip -9 > vars.$c.vcf.gz ; done
wget -O - ftp://ftp-trace.ncbi.nih.gov/1000genomes/ftp/release/20130502/ALL.chrX.phase3_shapeit2_mvncall_integrated_v1b.20130502.genotypes.vcf.gz | zcat - | cut -f1,2,3,4,5,6,7,8,1762 | gzip -9 > vars.X.vcf.gz
wget -O - ftp://ftp-trace.ncbi.nih.gov/1000genomes/ftp/release/20130502/ALL.chrY.phase3_integrated_v1b.20130502.genotypes.vcf.gz | zcat - | cut -f1,2,3,4,5,6,7,8,1762 | gzip -9 > vars.Y.vcf.gz
	
	# TODO: Run full analysis on cloud without needing to download these to local.


# download reads (todo)

OUTDIR=/Users/cb/code/src/github.com/cb01/lxy/data/GM12878/hic/test
for A in SRR927086; do fastq-dump -X 5 --accession $A --gzip --outdir $OUTDIR; done


# Download error corrected PacBio reads
ftp://ftp-trace.ncbi.nih.gov/giab/ftp/data/NA12878/NA12878_PacBio_MtSinai/corrected_reads_gt4kb.fasta


