import { useEffect, useState } from 'react';

export interface PostSummary {
  id: string;
  title: string;
  createdAt: string;
}

export default function fetchAll(
  url: string,
  from = '2000-01-01',
  to = '2030-10-01',
  page = 0,
  pageSize = 2000,
): [PostSummary[], boolean] {
  const [posts, setPosts] = useState([]);
  const [isSyncing, setIsSyncing] = useState(true);
  fetch(`${url}?from=${from}&to=${to}&page=${page}&page_size=${pageSize}`)
  .then(async (res: Response) => {
    setPosts(await res.json());
  })
  .catch(e => console.error(e))
  .finally(() => {
    setIsSyncing(false);
  });
  return [posts, isSyncing];
}
