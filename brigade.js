const { events, Job, Group } = require("brigadier");

events.on("check_suite:requested", runSuite);
events.on("check_suite:rerequested", runSuite);
events.on("check_run:rerequested", runSuite);

function runSuite(e, p) {
    // Run four tests in parallel. Each will report its results to GitHub.
    runDCO(e, p).catch(e => {console.error(e.toString())});
    runUnitTests(e, p).catch(e => {console.error(e.toString())});
    runStyleTests(e, p).catch(e => {console.error(e.toString())});
    runCoverage(e, p).catch(e => {console.error(e.toString())});
}

// Not as cool as runDMC.
function runDCO(e, p) {
    const re = new RegExp(/Signed-off-by:(.*)/, 'i');
    const ghData = JSON.parse(e.payload);
    const commitMsg = ghData.body.check_suite.head_commit.message
    const signedOff = re.exec(commitMsg);

    var dco = new Notification("dco", e, p);
    dco.title = "Developer Certificate of Origin (DCO)"

    if (signedOff == null){
        dco.summary = "DCO check failed: Not signed off."
        dco.text = "This commit is inelligible for merging until it is signed off. See https://developercertificate.org/\n\n```"+commitMsg+"```";
        dco.conclusion = "failure";
    } else {
        dco.summary = "DCO check succeeded";
        dco.text = `Commit signed by ${ signedOff[1] }`;
        dco.conclusion = "success";
    }
    
    console.log(dco.summary);
    return dco.run();
}

// Run tests and fail if the tests do not pass.
function runUnitTests(e, p) {
    const command = "make test";
    var note = new Notification("tests", e, p);
    note.conclusion = "";
    note.title = "Run Tests"
    note.summary = "Running the test target for " + e.revision.commit;
    note.text = "This test will execute all of the unit tests for the program."

    var job = new GoJob("run-tests", e, p);
    job.tasks.push(command);
    return notificationWrap(job, note)
}

// Run style tests
async function runStyleTests(e, p) {
    // Thank you Radu! 
    const commands = [
        "go get honnef.co/go/tools/cmd/gosimple",
        "go get -u golang.org/x/tools/cmd/gotype",
        "go get github.com/fzipp/gocyclo",
        "go get github.com/gordonklaus/ineffassign",
        "go get honnef.co/go/tools/cmd/unused",
        "go get -u github.com/kisielk/errcheck",
        "go get -u github.com/mibk/dupl",
        "go get honnef.co/go/tools/cmd/staticcheck",
        "go get github.com/walle/lll/...",
        "set +e",
        "export pkg=$(go list .)",
        `go vet "$pkg";
        gotype $GOPATH/src/"$pkg";
        gocyclo -over 10 $GOPATH/src/"$pkg";
        gosimple "$pkg";
        ineffassign $GOPATH/src/"$pkg";
        unused "$pkg";
        errcheck "$pkg";
        dupl $GOPATH/src/"$pkg";
        staticcheck "$pkg";
        lll --maxlength 120  $GOPATH/src/"$pkg";
        exit 0`
    ];
    
    var note = new Notification("style", e, p);
    note.conclusion = "";
    note.title = "Run Style Tests"
    note.summary = "Running the style test target for " + e.revision.commit;
    note.text = "This test checks for formatting, dead code, and other frequent problems."
    var job = new GoJob("run-style", e, p);
    Array.prototype.push.apply(job.tasks, commands);

    return notificationWrap(job, note);
}

// Test code coverage
function runCoverage(e, p) {
    var note = new Notification("coverage", e, p);
    note.conclusion = "";
    note.title = "Run Coverage Check"
    note.summary = "Running the test coverage report for " + e.revision.commit;
    note.text = "This test checks to see how much of the code has test coverage."
    var job = new GoJob("run-coverage", e, p);
    job.tasks.push("go test -cover .");

    return notificationWrap(job, note, "neutral");
}

// Helper to wrap a job execution between two notifications.
async function notificationWrap(job, note, conclusion) {
    if (conclusion == null) {
        conclusion = "success"
    }
    await note.run();
    try {
        let res = await job.run()
        const logs = await job.logs();

        note.conclusion = conclusion;
        note.summary = `Task "${ job.name }" passed`;
        note.text = note.text = "```" + res.toString() + "```\nTest Complete";
        return await note.run();
    } catch (e) {
        const logs = await job.logs();
        note.conclusion = "failure";
        note.summary = `Task "${ job.name }" failed for ${ e.buildID }`;
        note.text = "```" + logs + "```\nFailed with error: " + e.toString();
        try {
            return await note.run();
        } catch (e2) {
            console.error("failed to send notification: " + e2.toString());
            console.error("original error: " + e.toString());
            return e2;
        }
    }
}

// A GitHub Check Suite notification
class Notification {
    constructor(name, e, p) {
        this.proj = p;
        this.payload = e.payload;
        this.name = name;
        this.externalID = e.buildID;
        this.detailsURL = `https://azure.github.com/kashti/builds/${ e.buildID }`;
        this.title = "running check";
        this.text = "";
        this.summary = "";

        // count allows us to send the notification multiple times, with a distinct pod name
        // each time.
        this.count = 0;

        // One of: "success", "failure", "neutral", "cancelled", or "timed_out".
        this.conclusion = "neutral";
    }

    // Send a new notification, and return a Promise<result>.
    run() {
        this.count++
        var j = new Job(`${ this.name }-${ this.count }`, "technosophos/brigade-github-check-run:latest");
        j.env = {
            CHECK_CONCLUSION: this.conclusion,
            CHECK_NAME: this.name,
            CHECK_TITLE: this.title,
            CHECK_PAYLOAD: this.payload,
            CHECK_SUMMARY: this.summary,
            CHECK_TEXT: this.text,
            CHECK_DETAILS_URL: this.detailsURL,
            CHECK_EXTERNAL_ID: this.externalID
        }
        return j.run();
    }
}

// A Go-based Job
class GoJob extends Job {
    constructor(name, e, project) {
        super(name, "golang:1.9");

        this.e = e;
        this.project = project;
        const gopath = "/go"
        const localPath = gopath + "/src/github.com/" + project.repo.name;
        this.tasks = [
            "curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh",
            "mkdir -p " + localPath,
            "mv /src/* " + localPath,
            "cd " + localPath,
            "dep ensure",
        ];
    }
}