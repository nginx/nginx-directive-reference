'use client'
import { getDirectives } from '@nginx/reference-lib'
import { useCallback, useState } from 'react'
import { SearchBox } from './SearchBox'
import { Link } from 'react-router-dom'
import { toHref } from './params'

const directivesByModule = getDirectives().reduce((acc, d) => {
  const k = d.module
  if (acc.has(k)) {
    acc.get(k)?.push(d.name.toLowerCase())
  } else {
    acc.set(k, [d.name.toLowerCase()])
  }
  return acc
}, new Map<string, string[]>())

// alphabetize the module/directive tree
const modules = Array.from(directivesByModule.entries())
  .sort((a, b) => (a[0] > b[0] ? 1 : -1)) // sort by key
  .map(([name, directives]) => {
    return {
      id: name,
      name,
      directives: directives.sort(),
    }
  })

type ModList = typeof modules

type DirectiveNodeProps = {
  name: string
  modId: string
}

function DirectiveNode({ name, modId }: DirectiveNodeProps) {
  return (
    <li>
      <Link to={toHref({ modId, directive: name })}>
        <code>{name}</code>
      </Link>
    </li>
  )
}

type ModNodeProps = {
  id: string
  name: string
  open: boolean
  directives: string[]
  onClick: (id: string, directive?: string) => void
}

function ModNode({ id, name, directives, open, onClick }: ModNodeProps) {
  return (
    <li>
      <span
        className={`caret ${open ? 'open' : ''}`}
        onClick={() => onClick(id)}
      >
        <code>{name}</code>({directives.length ?? 0})
      </span>
      {open && (
        <ul>
          {directives.map((d) => (
            <DirectiveNode key={d} name={d} modId={id} />
          ))}
        </ul>
      )}
    </li>
  )
}

type TreeProps = {
  modId?: string
  onClick?: (modId: string, directive: string) => void
}

/** Tree shows an expandable tree of modules and directives */
export default function Tree({ modId, onClick }: TreeProps) {
  const [openModIds, setOpenModIds] = useState(modId ? [modId] : [])
  const [mods, setMods] = useState(modules)

  const isOpen = (id: string) => openModIds.includes(id)
  const toggleOpen = (id: string) => {
    const newModIds = isOpen(id)
      ? openModIds.filter((x) => x !== id)
      : [id, ...openModIds]
    setOpenModIds(newModIds)
  }
  const handleClick = (modId: string, directive?: string) => {
    if (directive && onClick) {
      onClick(modId, directive)
    }
    {
      toggleOpen(modId)
    }
  }

  // useCallback needed to prevent an loop between SearchBox's useEffect and onSearch
  const onSearch = useCallback(
    (term: string) => {
      if (term === '') {
        setMods(modules)
        setOpenModIds([])
        return
      }

      const filtered = modules.reduce((acc, m) => {
        const matches = m.directives.filter((d) =>
          d.includes(term.toLowerCase())
        )
        if (matches.length === 0) {
          return acc
        }
        const filteredMod = { directives: matches, id: m.id, name: m.name }
        return [...acc, filteredMod]
      }, [] as ModList)

      setMods(filtered)
      setOpenModIds(filtered.map((m) => m.id))
    },
    [setMods, setOpenModIds]
  )

  return (
    <div>
      <SearchBox onChange={onSearch} />
      {mods.length === 0 ? (
        <p>No matches found</p>
      ) : (
        <ul className="modules">
          {mods.map((m) => (
            <ModNode
              key={m.id}
              {...m}
              open={isOpen(m.id)}
              onClick={handleClick}
            />
          ))}
        </ul>
      )}
    </div>
  )
}
