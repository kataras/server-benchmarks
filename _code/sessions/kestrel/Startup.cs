using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Hosting;
using Microsoft.AspNetCore.Routing;

namespace netcore
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        public void ConfigureServices(IServiceCollection services)
        {
            services.AddRouting();
            // Adds a default in-memory implementation of IDistributedCache.
            services.AddDistributedMemoryCache();

            services.AddSession(options =>
            {
                options.Cookie.Name = "session";
            });
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            var routeBuilder = new RouteBuilder(app);

            routeBuilder.MapGet("/sessions", context =>{
                context.Session.SetString("ID", System.Guid.NewGuid().ToString());

                context.Session.SetString("name", "John Doe");
                var name = context.Session.GetString("name");
                return context.Response.WriteAsync(name);
            });

            var routes = routeBuilder.Build();

            app.UseSession();
            app.UseRouter(routes);
        }
    }
}
