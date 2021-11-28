import config from '~/lib/config.ts'
import fetchAll from '~/lib/fetchAll.ts'
import Footer from '~/components/footer.tsx'
import MainTop from '~/components/mainTop.tsx'
import React from 'react'
import SummaryList from '~/components/summaryList.tsx'

export default function Home() {
  const envs = config();
  const [postsByDate, isSyncing] = fetchAll(
    `${envs.TRANSPORT}://${envs.GATEWAY_HOST}:${envs.GATEWAY_PORT}/posts-by-date`,
  );

  return (
    <div className="page">
      <MainTop logoCls="logo" linksCls="links"/>
      <div className="entries">
        {isSyncing && (
          <p className="entries-txt">...</p>
        )}
        {!isSyncing && postsByDate.length > 0 && (
          SummaryList("list", postsByDate)
        )}
        {!isSyncing && postsByDate.length === 0 && (
         <p className="entries-txt">Nothing to see here, for now...</p>
        )}
      </div>
      <Footer />
    </div>
  )
}
