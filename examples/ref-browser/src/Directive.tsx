import { Format, getDirectives } from '@nginx/reference-lib'
import Description from './Description'

type Dir = {
  name: string
  contexts: string[]
  syntax: string[]
  default?: string
}
type DirIndex = Record<string, Dir>
type ModIndex = Record<string, DirIndex>

// build an index of modules and directives
const index = getDirectives(Format.HTML).reduce((mods, m) => {
  if (m.module in mods) {
    mods[m.module][m.name] = m
  } else {
    mods[m.module] = {
      [m.name]: m,
    }
  }
  return mods
}, {} as ModIndex)

type SyntaxProps = {
  directive: string
  syntax: string
}
function Syntax({ directive, syntax }: SyntaxProps) {
  return (
    <div className="syntax">
      <code>{directive}</code>
      {syntax && <span dangerouslySetInnerHTML={{ __html: syntax }} />}
      <code>;</code>
    </div>
  )
}

type DirectiveProps = {
  modId: string
  directive: string
}
export default function Directive({ modId, directive }: DirectiveProps) {
  if (!(modId in index)) {
    return (
      <p>
        No such module <code>{modId}</code>
      </p>
    )
  }
  const dir = index[modId]?.[directive]
  if (!dir) {
    return (
      <p>
        No such directive <code>{directive}</code>
      </p>
    )
  }

  return (
    <div style={{ width: '100%' }}>
      <h2>
        <code>{dir.name}</code>
      </h2>
      <dl>
        <dt title={JSON.stringify(dir.syntax)}>Syntax</dt>
        <dd>
          {dir.syntax.map((s, i) => (
            <Syntax key={i} directive={dir.name} syntax={s} />
          ))}
        </dd>
        {dir.default && (
          <>
            <dt title={JSON.stringify(dir.default)}>Default</dt>
            <dd>{dir.default}</dd>
          </>
        )}
        {dir.contexts && (
          <>
            <dt title={JSON.stringify(dir.contexts)}>Contexts</dt>
            <dd>
              {dir.contexts.map((c) => (
                <code className="context" key={c}>
                  {c}
                </code>
              ))}
            </dd>
          </>
        )}
      </dl>
      <Description module={modId} directive={directive} />
    </div>
  )
}
