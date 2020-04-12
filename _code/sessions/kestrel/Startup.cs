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
        private static byte[] BodyBytes = System.Text.Encoding.UTF8.GetBytes("John Doe");

        public SessionMiddleware(RequestDelegate next)
        {
        }

        public Task Invoke(HttpContext context)
        {
            context.Session.SetString("ID", context.TraceIdentifier);

            context.Session.Set("name", BodyBytes);
            context.Session.TryGetValue("name", out var name);
            return context.Response.Body.WriteAsync(name, 0, name.Length);
        }
    }
}
