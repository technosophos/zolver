// Acid CI/CD
console.log("===> Building " + pushRecord.repository.full_url + " " + pushRecord.head_commit.id);

// We need to simulate a Go environment
gopath = "/go";
localPath = gopath + "/src/github.com/" + pushRecord.repository.full_name;

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

j.run(pushRecord).waitUntilDone()
