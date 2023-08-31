import { find, Format, getDirectives, Directive} from '../index';
import mockReference from '../src/__mocks__/reference_mock.json'

describe('Directive Helper', () => {
  const directive = mockReference.modules.at(0)?.directives.at(0)
  describe('find', () => {
    test('with module', () => {
      const actual = find('allow', 'ngx_http_access_module', Format.HTML)
      const expected = directive?.description_html
      expect(actual).toBe(expected)
    });
    test('without module', () => {
      const actual = find('allow', undefined, Format.HTML)
      const expected = directive?.description_html
      expect(actual).toBe(expected)
    });
    test('returns undefined if not found', () => {
      const actual = find('listen', '', Format.HTML)
      expect(actual).toBeUndefined
    });
    test('returns HTML', () => {
      const actual = find('allow', '', Format.HTML)
      const expected = directive?.description_html
      expect(actual).toBe(expected)
    });
    test('returns Markdown', () => {
      const actual = find('allow', '', Format.Markdown)
      const expected = directive?.description_md
      expect(actual).toBe(expected)
    });
  })

  describe('getDirectives', () => {
    const module = mockReference.modules.at(0)
    const directive = module?.directives.at(0)
    test('returns HTML', () => {
      const actual = getDirectives(Format.HTML)
      const expected = [{ name: directive?.name,
        module: module?.name,
        description: directive?.description_html,
        syntax: directive?.syntax_html,
        contexts: directive?.contexts,
        isBlock: directive?.isBlock,
        default: directive?.default,
      } as Directive ]
      expect(actual).toStrictEqual(expected)
    });
    test('returns Markdown', () => {
      const actual = getDirectives(Format.Markdown)
      const expected = [{ name: directive?.name,
        module: module?.name,
        description: directive?.description_md,
        syntax: directive?.syntax_md,
        contexts: directive?.contexts,
        isBlock: directive?.isBlock,
        default: directive?.default,
        } as Directive ]
      expect(actual).toStrictEqual(expected)
    });
  })
})
