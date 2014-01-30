package main

import (
	"flag"
	"fmt"
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os/exec"
)

var (
	port    = flag.Int("p", 8080, "webserver port")
	dirPath = flag.String("dir", "./up/", "directory path for uploaded files")
)

type Server struct {
	m       *martini.ClassicMartini
	session *mgo.Session
	db      *mgo.Database
	gfs     *mgo.GridFS
}

type Result struct {
	Files []*FileInfo `json:"files"`
}

type FileInfo struct {
	Url          string `json:"url,omitempty"`
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Size         int64  `json:"size"`
	Error        string `json:"error,omitempty"`
	DeleteUrl    string `json:"delete_url,omitempty"`
	DeleteType   string `json:"delete_type,omitempty"`
}

func NewServer() *Server {
	var err error
	var s Server
	s.session, err = mgo.Dial("mongodb://admin:admin@localhost/test")
	if err != nil {
		fmt.Println("Cant connect to Mongo")
		panic(err)
	}
	s.db = s.session.DB("")
	s.gfs = s.db.GridFS("fs")

	s.m = martini.Classic()

	s.m.Handlers(martini.Recovery(), martini.Static("public"))

	//s.m.Get("/", func() {
	//	println(indexPage)
	//})
	s.registerHandlers()
	//s.m.Post("/filename/:filename", printFile)
	return &s
}
func (s *Server) registerHandlers() {
	s.m.Get("/", s.Handler)
	s.m.Post("/printer/:printername/file/:filename", s.PrintFile)
}

func (s *Server) Run() {
	log.Print("api: starting server")

	s.m.Run()
}

func main() {
	flag.Parse()
	if *dirPath == "" {
		log.Fatal("Please specify directory path for uploaded files")
	}

	var s *Server = NewServer()
	s.Run()
	//For listing the available GridFS objects

	defer s.session.Close()

	var result *mgo.GridFile
	iter := s.gfs.Find(nil).Iter()

	for s.gfs.OpenNext(iter, &result) {
		fmt.Printf("Filename: %s\n", result.Name())
	}
	if iter.Err() != nil {
		panic(iter.Err())
	}
	//

	//http.HandleFunc("/", handleHome)
	/*addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Failed to run server: ", err)
	}*/

}

func (s *Server) Handler() string {

	//	fmt.Fprintf("%s", indexPage)
	return indexPage
}

func (s *Server) PrintFile(w http.ResponseWriter, params martini.Params) {
	fileName, printerName := params["filename"], params["printername"]

	PrintFile := exec.Command("lp", "-d", printerName, fileName)

	err := PrintFile.Run()
	if err != nil {

		fmt.Println("Error!")
	}
}

/*func handleHome(w http.ResponseWriter, r *http.Request) {
m.Get("/", fmt.Fprintf(w, indexPage))

 else {
	var result Result
	mr, _ := r.MultipartReader()
	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}
		if name := part.FormName(); name != "" {
			if part.FileName() != "" {
				result.Files = append(result.Files, uploadFile(w, part))
			}
		}
	}

	js, _ := json.Marshal(result)
	jsonType := "application/json"
	if strings.Index(r.Header.Get("Accept"), jsonType) != -1 {
		w.Header().Set("Content-Type", jsonType)
	}
	fmt.Fprintln(w, string(js))
} */

/*func uploadFile(w http.ResponseWriter, p *multipart.Part) (fi *FileInfo) {
	filePath := filepath.Join(*dirPath, p.FileName())
	fi = &FileInfo{
		Name: p.FileName(),
		Type: p.Header.Get("Content-Type"),
	}

	gfsFile, err := s.gfs.Create(p.FileName())

	size, err := io.Copy(gfsFile, p)
	if err == nil {
		fi.Size = size
		fmt.Println("Uploaded to GridFS")
	} else {
		fmt.Println("Failed....")
	}

	if f, e := os.Create(filePath); e == nil {
		size, err := io.Copy(f, p)
		if err == nil {
			fi.Size = size
			fmt.Println("Uploaded")
		} else {
			fmt.Println("Failed....")
		}
	}
	return
}
*/
const indexPage = `
<!DOCTYPE HTML>
<html>
<head>
<meta charset="utf-8">
<title>CCLOM File Uploader</title>
<script src="//ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
<script src="//cdn.jsdelivr.net/jquery.fileupload/8.9.0/js/vendor/jquery.ui.widget.js"></script>
<script src="//cdn.jsdelivr.net/jquery.fileupload/8.9.0/js/jquery.iframe-transport.js"></script>
<script src="//cdn.jsdelivr.net/jquery.fileupload/8.9.0/js/jquery.fileupload.js"></script>
<script src="//cdn.jsdelivr.net/jquery.fileupload/8.9.0/js/jquery.fileupload-process.js"></script>
<link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.0.0/css/bootstrap.min.css">
<link rel="stylesheet" href="//cdn.jsdelivr.net/jquery.fileupload/8.9.0/css/jquery.fileupload.css">
<style>
body{
	padding:10px;
}
.bar {
    height: 18px;
    background: green;
}
</style>
</head>
<body>

<form id="fileupload" method="POST" enctype="multipart/form-data">
        <div class="row fileupload-buttonbar">
            <div class="col-lg-7">
                <!-- The fileinput-button span is used to style the file input field as button -->
                <span class="btn btn-success fileinput-button">
                    <i class="glyphicon glyphicon-plus"></i>
                    <span>Add files...</span>
                    <input type="file" name="files[]" multiple>
                </span>
                 <span class="btn btn-danger fileinput-button" id="refresh">
                    <i class="glyphicon glyphicon-plus"></i>
                    <span>Refresh</span>
                </span>
            </div>
          </div>
           <div class="row fileupload-buttonbar">
           <br/>
            <!-- The global progress state -->
            <div class="col-lg-4 fileupload-progress ">
                <!-- The global progress bar -->
                <div class="progress progress-striped active" role="progressbar" aria-valuemin="0" aria-valuemax="100">
                    <div class="progress-bar progress-bar-success" style="width:0%;"></div>
                </div>
                <!-- The extended global progress state -->
                <div class="progress-extended">&nbsp;</div>
            </div>
        </div>
        <!-- The table listing the files available for upload/download -->
        <table role="presentation" class="table table-striped"><tbody class="files"></tbody></table>
    </form>
<div class="row">
	<ul id='list'>
	</ul>
</div>
<script>
$(function () {
	
	$('#refresh').click(function(){
		location.reload();
	});

    $('#fileupload').fileupload({
        dataType: 'json',
        done: function (e, data) {
            $.each(data.result.files, function (index, file) {
                $('<li/>').text(file.name +" uploaded").appendTo($('#list'));
            });
        },
        progressall: function (e, data) {
	        var progress = parseInt(data.loaded / data.total * 100, 10);
	         if (progress){
	        	 var p = progress.toString()+"%%";
	        	 $('.progress-bar.progress-bar-success').width(p);
	    	}
   		 }
    });
});
</script>
</body> 
</html>
`
