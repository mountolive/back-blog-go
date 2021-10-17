import { useState } from 'react'
import transformDate from '~/lib/transformDate.ts'

export interface PostSummary {
  id: string;
  title: string;
  createdAt: string;
}

interface SummaryResponseDTO {
  id: string;
  title: string;
  created_at: string;
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

  if (isSyncing) {
    fetch(`${url}?from=${from}&to=${to}&page=${page}&page_size=${pageSize}`)
    .then((res: Response) => res.ok ? res.json() : Promise.reject(res))
    .then((allPosts: SummaryResponseDTO[]) => {
      setPosts(
        allPosts.map((post: SummaryResponseDTO) => {
          return {
            id: post.id,
            title: post.title,
            createdAt: transformDate(post.created_at),
          }
        }),
      )
    })
    .catch(e => console.error(e))
    .finally(() => {
      setIsSyncing(false);
    });
  }

  return [posts, isSyncing];
}
