using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.AspNetCore.Http;

namespace netcore
{
    public class Startup
    {
        public void ConfigureServices(IServiceCollection services)
        {
            // Adds a default in-memory implementation of IDistributedCache.
            services.AddDistributedMemoryCache();

            services.AddSession(options =>
            {
                options.Cookie.Name = "session";
            });
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            app.UseSession();
            app.UseMiddleware<SessionMiddleware>();
        }
    }

    public class SessionMiddleware
    {
        private static System.Text.Encoding Encoding = System.Text.Encoding.UTF8;

        private static readonly PathString SessionPath = new PathString("/session");

        private RequestDelegate _next;

        public SessionMiddleware(RequestDelegate next)
        {
            _next = next;
        }

        public Task Invoke(HttpContext context)
        {
            if (context.Request.Method != HttpMethods.Get || !SessionPath.StartsWithSegments("/session"))
                return _next(context);

            context.Session.Set("ID", Encoding.GetBytes(context.TraceIdentifier));
            context.Session.Set("name", Encoding.GetBytes("John Doe"));
            context.Session.TryGetValue("name", out var name);
            return context.Response.Body.WriteAsync(name, 0, name.Length);
        }
    }
}
