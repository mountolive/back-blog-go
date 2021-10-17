import React from 'react'

export default function GoBack() {
  return (
    <div>
      <p className="back" onClick={_event => location.href="/"}>Go back</p>
    </div>
  )
}
