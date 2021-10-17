
export default function transformDate(createdAt: string): string {
  const postDate = new Date(createdAt)
  const day = `${postDate.getDate() < 10 ? '0' : ''}${postDate.getDate()}`
  const month = `${(postDate.getMonth() + 1) < 10 ? '0' : ''}${postDate.getMonth() + 1}`

  return `${postDate.getFullYear()}-${month}-${day}`
}
