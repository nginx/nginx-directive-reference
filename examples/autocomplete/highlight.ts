import { getDirectives, Directive } from '@nginxinc/reference-lib'

/** kinds of tokens used for highlighting */
type token =
  | 'nginx.toplevel'
  | 'nginx.top.block'
  | 'nginx.block'
  | 'nginx.directives'

function toToken(d: Directive): token {
  if (d.contexts.includes('main')) {
    if (d.isBlock) {
      return 'nginx.top.block'
    }
    return 'nginx.toplevel'
  } else if (d.isBlock) {
    return 'nginx.block'
  }
  return 'nginx.directives'
}

function toRegex(directiveNames: Set<string>): RegExp {
  const escapedNames = Array.from(directiveNames)

  return new RegExp(`^\\s*(${escapedNames.join('|')})\\b`)
}

const dataset: Record<token, Set<string>> = {
  'nginx.toplevel': new Set<string>(),
  'nginx.top.block': new Set<string>(),
  'nginx.block': new Set<string>(),
  'nginx.directives': new Set<string>(),
}

for (const directive of getDirectives()) {
  const token = toToken(directive)
  dataset[token].add(directive.name)
}

const formatted: [RegExp, token][] = Object.entries(dataset).map(
  ([key, value]) => {
    return [toRegex(value), key as token]
  }
)

formatted.forEach(([regex, token]) => {
  console.log(`[${regex}, "${token}"],`)
})
