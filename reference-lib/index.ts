import reference from './src/reference.json'

export interface Directive {
    name: string
    module: string
    description: string
    syntax: string[]
    contexts: string[]
    isBlock: boolean
    default: string
}

export enum Format {
  Markdown = 1,
  HTML = 2,
}

type DirectiveHelp = {
  name: string
  description_md: string
  description_html: string
  module: string
}

// map that stores the directive name and array of helper content as
// directives names are not unique there can be multiple modules that share the same directive name
const refDirectives = new Map<string, DirectiveHelp[]>()
for (const modules of reference.modules) {
  if (modules.directives === null) continue
  for (const directive of modules.directives) {
    const data = refDirectives.get(directive.name) || []
    data.push({
      ...directive,
      module: modules.name,
    })
    refDirectives.set(directive.name, data)
  }
}

/**
 * Returns all the nginx directives
 *
 *  @param: format: format of the return type HTML or markdown
 *
 *  @return: an array of Directives
 */
export function getDirectives(format=Format.HTML): Directive[] {
    const directives = reference.modules.flatMap((m) =>
      m.directives.map((d) => ({...d, module: m.name})))
    .map ((d) => ({
        name: d.name,
        module: d.module,
        description: format === Format.HTML ? d.description_html : d.description_md,
        syntax: format === Format.HTML ? d.syntax_html : d.syntax_md,
        contexts: d.contexts,
        isBlock: d.isBlock,
        default: d.default
    } as Directive))
    return directives
}

/**
 * Returns the description corresponding to the directive name
 *
 *  @param: directive: directive name to find
 *  @param: module: optional name of module
 *  @param: format: format of the return type HTML or markdown
 *
 *  @return: a string containing the description of the directive in xml or markdown format
 */
export function find(directive: string, module: string | undefined, format=Format.HTML): string | undefined {
  const data =
    module
      ? refDirectives.get(directive)?.find((d) => d.module.toUpperCase() === module.toUpperCase())
      : refDirectives.get(directive)?.at(0)
  return (format === Format.HTML ? data?.description_html : data?.description_md)
}
