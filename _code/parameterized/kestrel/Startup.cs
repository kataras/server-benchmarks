using System;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Routing;

namespace netcore
{
    public class Startup
    {
        private static byte[] BodyBytes = System.Text.Encoding.UTF8.GetBytes("Hello ");

        public void Configure(IApplicationBuilder app)
        {
            var routeBuilder = new RouteBuilder(app).MapGet("/hello/{name}", context =>
            {
                var name = (string)context.GetRouteValue("name");
                var outputBytes = new byte[BodyBytes.Length + name.Length];
                Array.Copy(BodyBytes, outputBytes, BodyBytes.Length);
                var outputBytesCount = System.Text.Encoding.UTF8.GetBytes(name, outputBytes.AsSpan(BodyBytes.Length));

                return context.Response.Body.WriteAsync(outputBytes, 0, BodyBytes.Length + outputBytesCount);
            });

            app.UseRouter(routeBuilder.Build());
        }
    }
}