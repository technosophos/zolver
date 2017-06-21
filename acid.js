// Acid CI/CD
events.push = function(e) {
  console.log("===> Building " + e.repo.cloneURL + " " + e.commit);

  // We need to simulate a Go environment
  gopath = "/go";
  localPath = gopath + "/src/github.com/" + e.repo.name;

  // Define a single build step:
  j = new Job("test-zolver");
  j.image = "technosophos/acid-go:latest";
  j.env = {
    "DEST_PATH": localPath,
    "GOPATH": gopath
  };

  j.tasks = [
    "go get github.com/Masterminds/glide",
    "make bootstrap",
    "make test"
  ];

  j.run()
}
