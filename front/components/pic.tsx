import React from 'react'

export default function Pic({ size = 75, className }: { size?: number, className: string }) {
  return (
    <img src="/dechirico.jpeg" className={className} height={size} title="DeChirico" />
  )
}
