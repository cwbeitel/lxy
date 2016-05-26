import numpy as np
import numpy.random
import matplotlib.pyplot as plt
import sys
import click

# Visualize a dotplot to compare two orderings of a list
# e.g. python scripts/scaffplot.py --inferred data/test/scaffolding.inferred.txt --actual data/test/scaffolding.key.txt --outpath data/test/test2.png

@click.command()
@click.option('--inferred', required=True, help='A file containing a newline-delimited list of elements.')
@click.option('--actual', required=True, default="", help='A file containing a newline-delimited list of elements.')
@click.option('--label', default="unlabeled", help='Label to place in title of plot.')
@click.option('--outpath', required=True, default="", help='Path to write the output figure.')
def dotplot(inferred, actual, label, outpath):
    """Visualize a order comparison dotplot for two ordered lists of ids."""

    data = {}
    id_count = 0

    # Scan through the actual/key list
    with open(actual, "r") as f:
        for line in f:
            key = line.strip()
            assert key not in data, "Found duplicate id in actual/key list."
            data[key] = {}
            data[key]["actual"] = id_count
            id_count += 1

    id_count_max = id_count

    id_count = 0

    # Scan through to get count of unique ids
    with open(inferred, "r") as f:
        for line in f:
            key = line.strip()
            assert key in data, "Actual and inferred id lists contain different elements. Key does not include " + key
            data[key]["inferred"] = id_count
            id_count += 1

    id_count_max = max(id_count_max, id_count)

    x = np.zeros(id_count_max + 1)
    y = np.zeros(id_count_max + 1)

    for k, v in data.iteritems():
        if "inferred" in v:
            actual = v["actual"]
            print actual
            print v["inferred"]
            x[actual-1] = actual
            y[actual-1] = v["inferred"]

    plt.clf()
    plt.plot(x, y, "bo")
    plt.title('Order evaluation')

    frame = plt.gca()
    frame.axes.get_xaxis().set_ticks([])
    frame.axes.get_yaxis().set_ticks([])

    plt.xlabel("inferred")
    plt.ylabel("actual")

    plt.savefig(outpath)

if __name__ == '__main__':
    dotplot()
