using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Routing;

namespace netcore
{
    public class Startup
    {
        private static System.Text.Encoding Encoding = System.Text.Encoding.UTF8;

        public void Configure(IApplicationBuilder app)
        {
            var routeBuilder = new RouteBuilder(app).MapGet("/hello/{name}", context =>
            {
                var name = (string)context.GetRouteValue("name");

                var outputBytes = Encoding.GetBytes("Hello " + name);

                return context.Response.Body.WriteAsync(outputBytes, 0, outputBytes.Length);
            });

            app.UseRouter(routeBuilder.Build());
        }
    }
}