// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.865
package gen

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	VM "github.com/CP-Payne/wonderpicai/web/template/viewmodel"
)

func GenForm(data VM.GenFormComponentData) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<form id=\"gen-form\" hx-post=\"/gen\" hx-swap=\"outerHTML\" hx-indicator=\"#generation-spinner-area\" class=\"space-y-6\" data-cost-per-image=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", data.Form.MinCost))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 12, Col: 61}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\"><div><h2 class=\"text-xl font-semibold mb-1 text-secondary\">Generation Settings</h2><p class=\"text-xs text-base-content/70\">Configure your AI image creation.</p></div><div><div class=\"inline-flex stats shadow bg-base-200 stats-horizontal w-full rounded-xl\"><div class=\"stat p-3 sm:p-4\"><div class=\"stat-figure text-primary\"><i class=\"fa-solid fa-cubes text-2xl\"></i></div><div class=\"stat-title text-xs sm:text-sm\">Available Credits</div><div id=\"available-credits\" class=\"stat-value text-lg sm:text-2xl\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d",
			data.Form.Credits))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 26, Col: 38}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</div><div class=\"stat-desc text-xs\">Ready to use</div></div><div class=\"stat p-3 sm:p-4\"><div class=\"stat-figure text-secondary\"><i class=\"fa-solid fa-wand-magic-sparkles text-2xl\"></i></div><div class=\"stat-title text-xs sm:text-sm\">Cost Per Image</div><div class=\"stat-value text-lg sm:text-2xl\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", data.Form.MinCost))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 35, Col: 98}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</div><div class=\"stat-desc text-xs\"><i class=\"fa-solid fa-cubes text-primary text-xs\"></i> / image</div></div><div class=\"stat p-3 sm:p-4 min-w-[160px]\"><div class=\"stat-figure text-accent flex items-center justify-center min-w-[32px] sm:min-w-[48px]\"><i class=\"fa-solid fa-coins text-2xl text-accent shirnk-o\"></i></div><div class=\"stat-title text-xs sm:text-sm\">Estimated Total Cost</div><div id=\"estimated-cost\" class=\"stat-value text-lg sm:text-2xl\">0</div><div class=\"stat-desc text-xs\">For this generation</div></div></div><label class=\"label\"><span class=\"text-xs text-base-content/60\">Scroll right to see calculated total cost --> </span></label></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if err, ok := data.Errors["credits"]; ok {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<p class=\"text-error text-xs -mt-4 ml-1\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(err)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 55, Col: 50}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "<div class=\"form-control w-full\"><label class=\"label mb-2\" for=\"prompt-textarea\"><span class=\"text-base font-semibold text-base-content\">Your Creative Prompt</span></label> <textarea id=\"prompt-textarea\" name=\"prompt\" class=\"textarea textarea-lg w-full h-36 sm:h-40 resize-y ring-1 ring-base-300 focus:ring-2 focus:ring-primary bg-base-100 placeholder:text-base-content/50\" placeholder=\"e.g. A hyperrealistic portrait of a wise old owl wearing a steampunk monocle, intricate details, dramataic lighting\" required rows=\"5\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(data.Form.Prompt)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 67, Col: 47}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "</textarea> <label class=\"label mt-1\"><span class=\"text-xs text-base-content/60\">The more detailed your prompt, the better the result</span></label> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if err, ok := data.Errors["prompt"]; ok {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "<p class=\"text-error text-xs mt-1\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 string
			templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(err)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 73, Col: 47}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "</div><div class=\"form-control\"><label class=\"label mb-2\" for=\"image-count-input\"><span class=\"text-base font-semibold text-base-content\">Number of Images</span> <span class=\"tooltip tooltip-left sm:tooltip-top\" data-tip=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("Each image costs %d credits.",
			data.Form.MinCost))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 82, Col: 34}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "\"><i class=\"fa-solid fa-circle-info text-base-content/70\"></i></span></label><div class=\"join w-full\"><button type=\"button\" id=\"decrement-images\" class=\"btn join-item btn-outline btn-secondary\">-</button> <input id=\"image-count-input\" type=\"number\" name=\"image_count\" class=\"input join-item w-full text-center font-semibold ring-1 ring-base-300 focus:ring-2 focus:ring-primary bg-base-100\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var9 string
		templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", data.Form.ImageCount))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 90, Col: 63}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "\" min=\"1\" max=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d",
			data.Form.MaxImagesPerGen))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 91, Col: 42}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "\" required> <button type=\"button\" id=\"increment-images\" class=\"btn join-item btn-outline btn-secondary\">+</button></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if err, ok := data.Errors["imageCount"]; ok {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "<p class=\"text-error text-xs mt-1\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var11 string
			templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(err)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/components/gen/gen_form.templ`, Line: 95, Col: 48}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "</div><div class=\"space-y-4 pt-6\"><button type=\"submit\" class=\"btn btn-primary btn-block text-base\"><span id=\"generation-spinner-area\" class=\"mr-2\"></span> <svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-5 w-5 mr-2\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\" stroke-width=\"2\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z\"></path></svg> Generate Images</button> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if data.Form.HasFailedImages {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "<button type=\"button\" hx-delete=\"/gen/image/failed\" hx-target=\"#gallery\" hx-swap=\"innerHTML\" hx-confirm=\"Are you sure you want to clear all failed image placeholders?\" class=\"btn btn-block btn-outline btn-error text-sm\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-5 w-5 mr-2\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\" stroke-width=\"2\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16\"></path></svg> Clear Failed Images</button>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "</div></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
