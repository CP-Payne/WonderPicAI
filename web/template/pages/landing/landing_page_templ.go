// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.865
package landing

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/CP-Payne/wonderpicai/web/template"

var appScreenshots = []struct {
	Src         string
	Alt         string
	Caption     string
	FeatureName string
}{
	{
		Src:         "/static/assets/images/screenshots/gen-page.png",
		Alt:         "WonderPicAI Image Generation Interface",
		Caption:     "Craft stunning visuals with our intuitive AI image generator. Describe your vision, and watch it come to life.",
		FeatureName: "Powerful AI Generation",
	},
	{
		Src:         "/static/assets/images/screenshots/credits-page.png",
		Alt:         "WonderPicAI Credits Purchase Page",
		Caption:     "Flexible credit packs to fuel your creativity. Get started easily and scale as you need.",
		FeatureName: "Simple Credit System",
	},
	{
		Src:         "/static/assets/images/screenshots/stripe-checkout.png",
		Alt:         "WonderPicAI Stripe Checkout Page",
		Caption:     "Secure payments powered by Stripe. Choose the plan that fits your creative needs and check out with confidence.",
		FeatureName: "Secure Stripe Checkout",
	},
}

func LandingPage() templ.Component {
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
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"hero min-h-[calc(100vh-var(--navbar-height,4rem))] bg-base-200 relative overflow-hidden\"><div class=\"hero-overlay bg-opacity-30 bg-gradient-to-br from-primary via-secondary to-accent mix-blend-multiply\"></div><div class=\"hero-content text-center text-neutral-content relative z-10 px-4\"><div class=\"max-w-2xl\"><h1 class=\"mb-5 text-5xl sm:text-6xl lg:text-7xl font-bold leading-tight\">Unleash Your Visual Imagination with <span class=\"text-transparent bg-clip-text bg-gradient-to-r from-secondary to-accent\">WonderPicAI</span></h1><p class=\"mb-8 text-lg sm:text-xl opacity-90\">Transform your ideas into breathtaking images with the power of artificial intelligence. Simple, fast, and endlessly creative.</p><a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 templ.SafeURL = templ.URL("/auth/signup")
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var3)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\" class=\"btn btn-primary btn-lg shadow-lg hover:scale-105 transform transition-transform duration-200\">Get Started For Free <svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-5 w-5 ml-2\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\" stroke-width=\"2\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M17 8l4 4m0 0l-4 4m4-4H3\"></path></svg></a><p class=\"text-xs mt-4 opacity-70\">No credit card required to start.</p></div></div><div class=\"absolute -bottom-1/4 -left-1/4 w-1/2 h-1/2 bg-primary/10 rounded-full blur-3xl opacity-50\"></div><div class=\"absolute -top-1/4 -right-1/4 w-1/2 h-1/2 bg-accent/10 rounded-full blur-3xl opacity-50\"></div></div><div id=\"features\" class=\"py-16 sm:py-24 bg-base-100 text-base-content\"><div class=\"container mx-auto px-4\"><div class=\"text-center mb-12 sm:mb-16\"><h2 class=\"text-3xl sm:text-4xl font-bold text-primary mb-3\">See WonderPicAI in Action</h2><p class=\"text-lg text-base-content/80 max-w-xl mx-auto\">Experience a seamless workflow from idea to incredible AI-generated art.</p></div><div class=\"space-y-16 sm:space-y-20\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for i, shot := range appScreenshots {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<div class=\"flex flex-col items-center gap-8 sm:gap-12 md:flex-row\"><div class=\"w-full md:w-1/2 lg:w-3/5\"><div class=\"rounded-xl shadow-2xl overflow-hidden border-4 border-base-300 transform transition-all duration-300 hover:scale-105 hover:shadow-primary/30\"><img src=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(shot.Src)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/pages/landing/landing_page.templ`, Line: 77, Col: 25}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\" alt=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var5 string
				templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(shot.Alt)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/pages/landing/landing_page.templ`, Line: 77, Col: 42}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "\" class=\"w-full h-auto object-contain img-crisp-hack\"></div></div><div class=\"w-full md:w-1/2 lg:w-2/5 text-center md:text-left\"><h3 class=\"text-2xl sm:text-3xl font-semibold text-secondary mb-4\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(shot.FeatureName)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/pages/landing/landing_page.templ`, Line: 81, Col: 90}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</h3><p class=\"text-base-content/90 leading-relaxed mb-6\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(shot.Caption)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/template/pages/landing/landing_page.templ`, Line: 83, Col: 20}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</p>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if i == 0 {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "<a href=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var8 templ.SafeURL = templ.URL("/gen")
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var8)))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "\" class=\"btn btn-secondary btn-outline\">Try the Generator <svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 ml-2\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\" stroke-width=\"2\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M13 10V3L4 14h7v7l9-11h-7z\"></path></svg></a>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "</div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "</div></div></div><div class=\"py-16 sm:py-24 bg-base-200 text-base-content\"><div class=\"container mx-auto px-4 text-center\"><h2 class=\"text-3xl sm:text-4xl font-bold text-primary mb-6\">Ready to Create Without Limits?</h2><p class=\"text-lg text-base-content/80 max-w-xl mx-auto mb-10\">Join thousands of creators and start generating your unique images today. Explore our flexible credit options.</p><div class=\"space-y-4 sm:space-y-0 sm:flex sm:justify-center sm:gap-4\"><a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var9 templ.SafeURL = templ.URL("/auth/signup")
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var9)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "\" class=\"btn btn-primary btn-lg\">Sign Up Now</a></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = template.Base(true).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
