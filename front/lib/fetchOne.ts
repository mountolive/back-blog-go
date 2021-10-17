import { useState } from 'react';
import transformDate from '~/lib/transformDate.ts'

interface PostResponseDTO {
  id: string,
  creator: string,
  title: string,
  content: string,
  tags: string[],
  created_at: string,
}

export interface Post {
  title: string,
  content: string,
  createdAt: string,
}

export default function fetchOne(url: string, id: string): [Post, boolean] {
  const [post, setPost] = useState({})
  const [isSyncing, setIsSyncing] = useState(true)
  
  if (isSyncing) {
    fetch(`${url}/${id}`)
    .then((res: Response) => res.ok ? res.json() : Promise.reject(res))
    .then(({ title, content, created_at }: PostResponseDTO) => {
      setPost({
        title,
        content,
        createdAt: transformDate(created_at),
      })
    })
    .catch((err: any) => console.error(err))
    .finally(() => {
      setIsSyncing(false)
    })
  }

  return [post, isSyncing]
}
