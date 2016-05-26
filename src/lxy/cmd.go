package main 

import (
    "os"
    
    "github.com/codegangsta/cli"
    //log "github.com/Sirupsen/logrus"
    util "sequtil"
    
    "lxy/scaff"
    "lxy/phase"
    //"lxy/cluster"
)



/*

Command structure:

lxy

    phase
            prep
            infer
            eval
            viz
            all

    scaff
            prep
            infer
            eval
            viz
            all

    vars
            ...

    seq (fasta, sam, fastq)
            stats
            mask
            partition
            trim
            filter

    reproduce
            beitel16

*/

func Control() {
    app := cli.NewApp()
    app.Name = "lxy"
    app.Usage = "Genome analysis."
    app.EnableBashCompletion = true
    app.Commands = []cli.Command{

        scaff.ScaffoldCommand(),
        phase.PhaseCommand(),
        //cluster.ClusterCommand(),
        //ReproduceCommand(),
        util.VarsCommand(),
        util.SeqCommand(),
    }
    app.Run(os.Args)
}



