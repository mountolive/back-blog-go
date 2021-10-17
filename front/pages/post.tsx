import MainTop from '~/components/mainTop.tsx'
import PostDetail from '~/components/postDetail.tsx'
import React from 'react'
import { useRouter } from 'aleph/react'

export default function Post() {
  const router = useRouter()
  const id: string = router.query.has('id') ? router.query.get('id') : ''

  return (
    <div className="page">
     <MainTop logoCls="logo" linksCls="links" />
     {id ? PostDetail("post", id) : (
       <>
         <p className="wrongtag">Something wrong happened :/</p>
         <p className="back" onClick={_event => location.href="/"}>Go back</p>
       </>
     )}
    </div>
  )
}
