package cluster

import (
	"github.com/codegangsta/cli"
	//"os"

	//"github.com/golang/glog"
)

/*
Command example:
*/

func ClusterCommand() cli.Command {
    return cli.Command{
		Name:  "cluster",
		Usage: "Cluster contigs into contact groups using Hi-C data",
		Subcommands: []cli.Command{
	        cli.Command{
				Name:  "markov",
				Usage: "Given a set of Hi-C link data, perform Markov Clustering.",
				Flags: []cli.Flag{
				  	cli.StringFlag{
				  		Name: "links", 
				  		Value: "", 
				  		Usage: "Path to the Hi-C links file.",
				  	},
				  	cli.StringFlag{
				  		Name: "output", 
				  		Value: "", 
				  		Usage: "Output path.",
				  	},
				},
				Action: markovClusterCommand,
    		},
	        cli.Command{
				Name:  "eval",
				Usage: "Evaluate a contig clustering.",
				Flags: []cli.Flag{
				  	cli.StringFlag{
				  		Name: "phasing", 
				  		Value: "", 
				  		Usage: "Path to phasing file.",
				  	},
				  	cli.StringFlag{
				  		Name: "key", 
				  		Value: "", 
				  		Usage: "Path to phasing key file.",
				  	},
				  	cli.StringFlag{
				  		Name: "output", 
				  		Value: "", 
				  		Usage: "Output path for evaluation stats file.",
				  	},
				},
				Action: evalClusterCommand,
    		},
    	},

    }

}

func markovClusterCommand(c *cli.Context) {


}

func evalClusterCommand(c *cli.Context) {


}

