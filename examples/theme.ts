import type { editor } from "monaco-editor/esm/vs/editor/editor.api";

export const themeConfig: editor.IStandaloneThemeData = {
    colors: {
        // 'attribute.value.unit': '#68217a'
    },
    base: "vs",
    inherit: true,
    rules: [
        {
            token: "nginx.toplevel",
            foreground: "#B58440",
            fontStyle: "bold",
        },
        {
            token: "nginx.top.block",
            foreground: "#4071B5",
            fontStyle: "bold",
        },
        {
            token: "nginx.block",
            foreground: "#B540AB",
            fontStyle: "bold",
        },
        {
            token: "nginx.directives",
            foreground: "#40b54a", 
            fontStyle: "bold",
        },
    ],
};

export const themeDarkConfig: editor.IStandaloneThemeData = {
    colors: {
        // 'attribute.value.unit': '#68217a'
    },
    base: "vs-dark",
    inherit: true,
    rules: [
        {
            token: "nginx.toplevel",
            foreground: "#B58440",
            fontStyle: "bold",
        },
        {
            token: "nginx.top.block",
            foreground: "#4071B5",
            fontStyle: "bold",
        },
        {
            token: "nginx.block",
            foreground: "#B540AB",
            fontStyle: "bold",
        },
        {
            token: "nginx.directives",
            foreground: "#40b54a", 
            fontStyle: "bold",
        },
    ],
};
