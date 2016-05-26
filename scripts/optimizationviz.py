import numpy as np
import numpy.random
import matplotlib.pyplot as plt
import sys
import click

# Visualize a dotplot to compare two orderings of a list
# e.g. python scripts/scaffplot.py --inferred data/test/scaffolding.inferred.txt --actual data/test/scaffolding.key.txt --outpath data/test/test2.png

@click.command()
@click.option('--scores', required=True, help='A file containing a newline-delimited list of elements.')
@click.option('--label', default="unlabeled", help='Label to place in title of plot.')
@click.option('--outpath', required=True, default="", help='Path to write the output figure.')
def optviz(scores, label, outpath):
    """Visualize a order comparison dotplot for two ordered lists of ids."""

    x=[]
    y=[]
    #y2=[]
    last = 0

    # Scan through the actual/key list
    with open(scores, "r") as f:
        for line in f:
            arr = line.strip().split(" ")
            if int(arr[0]) < int(last):
                break
            else:
                last = arr[0]
            x.append(int(arr[0]))
            y.append(float(arr[1]))
            #y2.append(float(arr[2]))

    data = {}
    id_count = 0

    plt.clf()
    plt.plot(x, y)
    plt.title('Optimization score progression, ' + label)

    frame = plt.gca()

    plt.xlabel("Iteration")
    plt.ylabel("Score")

    plt.savefig(outpath)

if __name__ == '__main__':
    optviz()
