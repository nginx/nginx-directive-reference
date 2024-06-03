import { useState } from 'react'
import { Format, find } from '@nginx/reference-lib'

type DescriptionProps = {
  module: string
  directive: string
}

export default function Description({ module, directive }: DescriptionProps) {
  const [format, setFormat] = useState(Format.HTML)
  const desc = find(directive, module, format)
  if (!desc) {
    throw new Error(`invalid mod/directive pair ${module}/${directive}`)
  }

  return (
    <div style={{ width: '100%' }}>
      <div role="tablist">
        <button
          role="tab"
          className={format === Format.HTML ? 'active' : ''}
          onClick={() => {
            setFormat(Format.HTML)
          }}
        >
          Description
        </button>
        <button
          role="tab"
          className={format === Format.Markdown ? 'active' : ''}
          onClick={() => {
            setFormat(Format.Markdown)
          }}
        >
          Markdown
        </button>
      </div>
      {format === Format.HTML && (
        <div role="tabpanel" dangerouslySetInnerHTML={{ __html: desc }} />
      )}
      {format === Format.Markdown && <pre role="tabpanel">{desc}</pre>}
    </div>
  )
}
