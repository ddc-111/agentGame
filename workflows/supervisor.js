export const meta = {
  name: "agentgame-supervisor",
  description: "Supervisor workflow: monitors main improvement loop, corrects failures, audits quality, proposes new requirements"
};

const BACKLOG_PATH = "IMPROVEMENT_BACKLOG.json";
const REPORT_PATH = "SELF_IMPROVE_REPORT.md";
const SUPERVISOR_LOG = "SUPERVISOR_LOG.md";
const MAIN_WORKFLOW_ID = "wf_130c58188ffePxykZ8KkXwTxXk";

async function loadBacklog() {
  const raw = await readFile(BACKLOG_PATH);
  if (!raw) return [];
  return JSON.parse(raw);
}

async function saveBacklog(backlog) {
  await writeFile(BACKLOG_PATH, JSON.stringify(backlog, null, 2));
}

async function appendSupervisorLog(entry) {
  let existing = await readFile(SUPERVISOR_LOG);
  if (!existing) {
    existing = "# Supervisor Log\n\nMonitoring workflow: " + MAIN_WORKFLOW_ID + "\n\n---\n\n";
  }
  existing += entry + "\n\n---\n\n";
  await writeFile(SUPERVISOR_LOG, existing);
}

async function scanCodebase() {
  const result = await agent(
    "Scan the agentGame codebase at D:\\dev\\my_project\\agentGame for issues. Check:\n\n1. Run: cd server && go build ./... (does it compile?)\n2. Run: cd server && go vet ./... (any vet warnings?)\n3. Run: cd server && go test ./internal/tests/ -v -count=1 -timeout 120s (test results)\n4. Check if SELF_IMPROVE_REPORT.md exists and read it for completed work\n5. Read IMPROVEMENT_BACKLOG.json for current task statuses\n6. Look at recent git log (git log --oneline -10) for recent commits\n7. Check for any uncommitted changes (git status)\n8. Scan server/internal/network/ for any files > 500 lines (code smell)\n9. Check if any Go files have obvious issues: unused imports, dead code, TODO comments\n\nReturn JSON with:\n- build_ok: boolean\n- vet_ok: boolean\n- test_results: {passed, failed, failures: [{name, error}]}\n- recent_commits: [{hash, message}]\n- uncommitted_files: [string]\n- large_files: [{path, lines}]\n- issues_found: [{severity, file, description}]\n- backlog_status: {pending, completed, failed}\n- recommendations: [string]",
    { subagent_type: "general", timeout_ms: 300000 }
  );
  return result;
}

async function proposeNewRequirements(scanResults) {
  const result = await agent(
    "You are a senior software architect reviewing the agentGame project.\n\nBased on these scan results:\n" + JSON.stringify(scanResults, null, 2) + "\n\nAnd the current backlog:\n" + (await readFile(BACKLOG_PATH) || "[]") + "\n\nPropose NEW improvement tasks that are NOT already in the backlog. Focus on:\n\n1. **Critical fixes**: build failures, test failures, security holes\n2. **Missing features**: what would make this framework competitive?\n3. **Code quality**: performance issues, error handling gaps, missing edge cases\n4. **Developer experience**: better tooling, documentation, debugging aids\n5. **Game framework features**: what would game developers expect?\n6. **AI-native features**: what makes this unique as an AI game framework?\n\nFor each proposed task, provide JSON:\n{\n  \"id\": \"IMP-XXX\",\n  \"title\": \"short title\",\n  \"priority\": number (1=highest),\n  \"category\": \"quality|feature|refactor|testing|security|devops|docs|performance\",\n  \"description\": \"detailed description\",\n  \"files\": [\"paths to modify\"],\n  \"test_file\": \"path or null\",\n  \"rationale\": \"why this matters\"\n}\n\nReturn an array of proposed tasks. Be specific and actionable. Only propose things not already in the backlog.",
    { subagent_type: "general", timeout_ms: 300000 }
  );
  return result;
}

async function fixFailedTask(task) {
  const result = await agent(
    "A task failed in the agentGame self-improvement loop. Investigate and fix it.\n\nFailed task:\n" + JSON.stringify(task, null, 2) + "\n\nProject at: D:\\dev\\my_project\\agentGame\n\nSteps:\n1. Read the SELF_IMPROVE_REPORT.md for details on what was attempted\n2. Read the files that were supposed to be modified\n3. Understand what went wrong (test failures, build errors, logic issues)\n4. Implement a corrected version\n5. Run: cd server && go build ./...\n6. Run: cd server && go test ./internal/tests/ -v -count=1 -timeout 120s\n7. If tests pass, run: cd server && go vet ./...\n\nReturn JSON with:\n- fixed: boolean\n- what_was_wrong: string\n- what_was_fixed: string\n- tests_pass: boolean\n- build_ok: boolean\n- files_modified: [string]",
    { subagent_type: "general", timeout_ms: 600000 }
  );
  return result;
}

async function commitIfReady(message) {
  const result = await agent(
    "Check if there are uncommitted changes in D:\\dev\\my_project\\agentGame and commit them.\n\nSteps:\n1. git status\n2. If changes exist: git add -u (plus new test files)\n3. git commit -m '" + message + "'\n4. git log --oneline -3\n\nReturn: {committed: boolean, hash: string, files: [string]}",
    { subagent_type: "general", timeout_ms: 120000 }
  );
  return result;
}

// ========== SUPERVISOR LOOP ==========

let iteration = 0;
const MAX_SUPERVISOR_CYCLES = 10;

await phase("Supervisor Starting");
await appendSupervisorLog("## Supervisor started at " + new Date().toISOString());

while (iteration < MAX_SUPERVISOR_CYCLES) {
  iteration++;
  await phase("Supervisor Cycle " + iteration);

  // Step 1: Scan codebase health
  await log("Cycle " + iteration + ": Scanning codebase...");
  const scanResult = await scanCodebase();
  let scan;
  try {
    scan = typeof scanResult === "string" ? JSON.parse(scanResult) : scanResult;
  } catch {
    scan = { raw: String(scanResult).substring(0, 2000), issues_found: [], recommendations: [] };
  }

  await appendSupervisorLog("## Cycle " + iteration + " - Scan\n```\n" + JSON.stringify(scan, null, 2).substring(0, 3000) + "\n```");

  // Step 2: Handle failed tasks
  const backlog = await loadBacklog();
  const failedTasks = backlog.filter(t => t.status === "failed");

  if (failedTasks.length > 0) {
    await log("Found " + failedTasks.length + " failed tasks. Attempting fixes...");
    for (const task of failedTasks) {
      await log("Fixing: " + task.id + " - " + task.title);
      const fixResult = await fixFailedTask(task);
      let fix;
      try {
        fix = typeof fixResult === "string" ? JSON.parse(fixResult) : fixResult;
      } catch {
        fix = { fixed: false, raw: String(fixResult).substring(0, 1000) };
      }

      if (fix.fixed) {
        task.status = "pending"; // Reset to pending for re-processing
        await appendSupervisorLog("### Fix: " + task.id + "\n- Fixed: YES\n- What was wrong: " + (fix.what_was_wrong || "unknown") + "\n- What was fixed: " + (fix.what_was_fixed || "unknown"));
      } else {
        await appendSupervisorLog("### Fix: " + task.id + "\n- Fixed: NO\n- Details: " + JSON.stringify(fix).substring(0, 500));
      }
    }
    await saveBacklog(backlog);
  }

  // Step 3: Check for stalled main workflow
  if (scan.backlog_status) {
    const { pending, completed, failed } = scan.backlog_status;
    await log("Backlog status - Pending: " + pending + ", Completed: " + completed + ", Failed: " + failed);

    if (pending === 0) {
      await log("All tasks processed! Proposing new requirements...");
    }
  }

  // Step 4: Propose new requirements
  await log("Proposing new requirements...");
  const newReqsResult = await proposeNewRequirements(scan);
  let newReqs;
  try {
    const parsed = typeof newReqsResult === "string" ? JSON.parse(newReqsResult) : newReqsResult;
    newReqs = Array.isArray(parsed) ? parsed : (parsed.tasks || parsed.proposals || [parsed]);
  } catch {
    newReqs = [];
    await appendSupervisorLog("### New requirements parse failed\nRaw: " + String(newReqsResult).substring(0, 1000));
  }

  if (newReqs.length > 0) {
    const currentBacklog = await loadBacklog();
    const existingIds = new Set(currentBacklog.map(t => t.id));
    let added = 0;

    for (const req of newReqs) {
      if (!req.id || existingIds.has(req.id)) {
        // Generate unique ID
        const maxNum = currentBacklog.reduce((max, t) => {
          const match = t.id && t.id.match(/IMP-(\d+)/);
          return match ? Math.max(max, parseInt(match[1])) : max;
        }, 0);
        req.id = "IMP-" + String(maxNum + 1).padStart(3, "0");
      }
      req.status = "pending";
      currentBacklog.push(req);
      existingIds.add(req.id);
      added++;
    }

    await saveBacklog(currentBacklog);
    await appendSupervisorLog("### New requirements added: " + added + "\n" + newReqs.map(r => "- " + r.id + ": " + (r.title || r.description || "unnamed")).join("\n"));
    await log("Added " + added + " new requirements to backlog");
  }

  // Step 5: Quality audit - check for issues scan found
  if (scan.issues_found && scan.issues_found.length > 0) {
    const criticalIssues = scan.issues_found.filter(i => i.severity === "critical" || i.severity === "high");
    if (criticalIssues.length > 0) {
      await log("Found " + criticalIssues.length + " critical issues! Adding to backlog...");
      const currentBacklog = await loadBacklog();
      const maxNum = currentBacklog.reduce((max, t) => {
        const match = t.id && t.id.match(/IMP-(\d+)/);
        return match ? Math.max(max, parseInt(match[1])) : max;
      }, 0);

      criticalIssues.forEach((issue, i) => {
        currentBacklog.push({
          id: "IMP-" + String(maxNum + 1 + i).padStart(3, "0"),
          title: "CRITICAL: " + (issue.description || issue.file || "unknown issue"),
          priority: 0,
          category: "critical",
          description: "Auto-detected by supervisor: " + JSON.stringify(issue),
          files: [issue.file || ""],
          test_file: null,
          status: "pending"
        });
      });

      await saveBacklog(currentBacklog);
    }
  }

  // Step 6: Commit any uncommitted supervisor work
  if (scan.uncommitted_files && scan.uncommitted_files.length > 0) {
    await log("Committing supervisor changes...");
    await commitIfReady("supervisor(cycle-" + iteration + "): backlog updates and fixes");
  }

  // Step 7: Report
  const updatedBacklog = await loadBacklog();
  const summary = {
    cycle: iteration,
    pending: updatedBacklog.filter(t => t.status === "pending").length,
    completed: updatedBacklog.filter(t => t.status === "completed").length,
    failed: updatedBacklog.filter(t => t.status === "failed").length,
    needs_fix: updatedBacklog.filter(t => t.status === "needs_fix").length,
    total: updatedBacklog.length,
    build_ok: scan.build_ok,
    test_passed: scan.test_results ? scan.test_results.passed : "?",
    test_failed: scan.test_results ? scan.test_results.failed : "?",
    new_proposals: newReqs.length,
    critical_issues: scan.issues_found ? scan.issues_found.filter(i => i.severity === "critical" || i.severity === "high").length : 0
  };

  await appendSupervisorLog("## Cycle " + iteration + " Summary\n```\n" + JSON.stringify(summary, null, 2) + "\n```");
  await log("Cycle " + iteration + " complete. Total backlog: " + summary.total + " tasks.");

  // If everything is done and no critical issues, we can rest
  if (summary.pending === 0 && summary.critical_issues === 0 && summary.failed === 0) {
    await log("All clear! No pending tasks, no critical issues. Supervisor can rest.");
    break;
  }

  // Brief pause before next cycle (workflow will naturally yield)
  await log("Waiting before next cycle...");
}

await phase("Supervisor Complete");
await appendSupervisorLog("## Supervisor completed " + iteration + " cycles");

return {
  cycles: iteration,
  finalStatus: "supervisor_complete"
};
