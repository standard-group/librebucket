export async function fetchCommitHistory(username, reponame) {
  const res = await fetch(`/api/v1/repos/${username}/${reponame}/commits`);
  if (!res.ok) throw new Error('Failed to fetch commit history');
  return res.json();
}

export async function fetchCommit(username, reponame, hash) {
  const res = await fetch(`/api/v1/repos/${username}/${reponame}/commits/${hash}`);
  if (!res.ok) throw new Error('Failed to fetch commit');
  return res.json();
}

export async function fetchCommitChanges(username, reponame, hash) {
  const res = await fetch(`/api/v1/repos/${username}/${reponame}/commits/${hash}/changes`);
  if (!res.ok) throw new Error('Failed to fetch commit changes');
  return res.json();
}

export async function createCommit(username, reponame, message, files, token) {
  const res = await fetch(`/api/v1/repos/${username}/${reponame}/commits`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ message, files })
  });
  if (!res.ok) throw new Error('Failed to create commit');
  return res.json();
}