import numpy as np
import numpy.random
import matplotlib.pyplot as plt
import sys
import click
import math

@click.command()
@click.option('--inpath', required=True, help='Path to the links file.')
@click.option('--outpath', required=True, default="", help='Path to write the output figure.')
@click.option('--label', default="unlabeled", help='Label to place on axes.')
@click.option('--log', default=True, help='whether to display in log scale.')
@click.option('--ordering', required=True, help='The order of contigs.')
def heatmap(inpath, outpath, label, log, ordering):
    """Visualize a heatmap for a three-column space delimited text file."""

    ids = {}
    id_count = 0

    # Scan through to get count of unique ids
    with open(ordering, "r") as f:
        for line in f:
            if line[0] is not "#":
                val = str(line).strip()
                if val not in ids:
                    ids[val] = id_count
                    id_count += 1

    x = np.zeros((id_count, id_count))

    # Scan through to load data into the array
    with open(inpath, "r") as f:
        for line in f:
            if line[0] != "#":
                arr = line.split(" ")
                if arr[0] != arr[1]:
                    ind1 = ids[arr[0]]
                    ind2 = ids[arr[1]]
                    if log:
                        x[ind1][ind2] = math.log(float(arr[2]))
                        x[ind2][ind1] = math.log(float(arr[2]))
                    else:
                        x[ind1][ind2] = float(arr[2])
                        x[ind2][ind1] = float(arr[2])
    plt.clf()
    plt.imshow(x)
    plt.title('Interaction Frequency')

    frame = plt.gca()
    frame.axes.get_xaxis().set_ticks([])
    frame.axes.get_yaxis().set_ticks([])

    plt.xlabel(label)
    plt.ylabel(label)

    plt.savefig(outpath, dpi = 400)

if __name__ == '__main__':
    heatmap()
