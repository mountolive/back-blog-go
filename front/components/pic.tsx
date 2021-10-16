import React from 'react'

export default function Pic({ size = 75, className }: { size?: number, className: string }) {
  return (
    <img src="/small_leo.png" className={className} height={size} title="leoponcio" />
  )
}
