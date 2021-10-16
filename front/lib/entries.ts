import {useEffect, useState} from 'react';

export interface PostSummary {
  id: string;
  title: string;
  createdAt: string;
}

interface PostResponseDTO {
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
    .then(async (res: Response) => {
      return res.ok ? res.json() : Promise.reject(res);
    })
    .then((allPosts: PostResponseDTO[]) => {
      setPosts(
        allPosts.map((post: PostResponseDTO) => {
          return {
            id: post.id,
            title: post.title,
            createdAt: new Date(post.created_at).toDateString()
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
