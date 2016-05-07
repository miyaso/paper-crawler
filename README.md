# paper-crawler

# Install
Assume that $GOPATH is set and $PATH includes $GOPATH/bin.

    $ go get github.com/miyaso/paper-crawler
    $ paper-crawler -h  

# Example
    $  paper-crawler get SpectralClustring 

# Usage
    NAME:
       crawl_google_scholar_client -
    
    USAGE:
       paper-crawler [global options] command [command options] [arguments...]
    
    VERSION:
       0.0.1

    COMMANDS:
        get	Crawling google scholar.
    
    GLOBAL OPTIONS:
       --dryrun, -d		dryrun
       --help, -h		show help
       --version, -v	print the version
    
    NAME:
       paper-crawler get - Crawling google scholar.
    
    USAGE:
       paper-crawler get [command options] [arguments...]
    
    OPTIONS:
       --maxRequestNum "10"		The maximum number of crawling artciles
       --maxRequestCitedNum "10"	The maximum number of articleId citing by crawling artciles
       --crawlInterval "3"		Interval seconds for crawling.
