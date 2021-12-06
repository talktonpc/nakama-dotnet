/**
* Copyright 2021 The Nakama Authors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

using System.Linq;
using Xunit;
using NakamaSync;
using System;
using System.Threading.Tasks;

namespace Nakama.Tests.Sync
{
    public class SyncTestHost
    {
        [Fact(Timeout = TestsUtil.TIMEOUT_MILLISECONDS)]
        private async void UniqueHostShouldBeChosen()
        {
            var testEnv = new SyncTestEnvironment(
                numClients: 2,
                creatorIndex: 0);

            await testEnv.StartAll();
            var env1 = testEnv.GetUserEnv(testEnv.GetCreatorPresence());
            var env2 = testEnv.GetUserEnv(testEnv.GetRandomNonCreatorPresence());

            Assert.True(env1.Match.IsSelfHost() || env2.Match.IsSelfHost());
            Assert.NotEqual(env1.Match.IsSelfHost(), env2.Match.IsSelfHost());

            testEnv.Dispose();
        }

        [Fact(Timeout = TestsUtil.TIMEOUT_MILLISECONDS)]
        private async void HostChangedDispatches()
        {
            var testEnv = new SyncTestEnvironment(
                numClients: 2,
                creatorIndex: 0);

            IMatch env1Match = await testEnv.StartCreate();
            var env1 = testEnv.GetUserEnv(testEnv.GetCreatorPresence());

            IHostChangedEvent hostChangedEvent = null;

            env1.Match.OnHostChanged += evt =>
            {
                hostChangedEvent = evt;
            };

            var env2 = testEnv.GetNonCreatorEnv();
            await env2.JoinMatch(env1Match.Id);

            await Task.Delay(1000);

            string oldhostId = string.Compare(env1.Self.UserId, env2.Self.UserId, StringComparison.InvariantCulture) > 0 ? env2.Self.UserId : env1.Self.UserId;
            string newHostId = string.Compare(env1.Self.UserId, env2.Self.UserId, StringComparison.InvariantCulture) > 0 ? env1.Self.UserId : env2.Self.UserId;
            // todo there is a race here...the hostchangedevent doesn't always dispatch before the continuation.
            Assert.NotNull(hostChangedEvent);
            Assert.Equal(hostChangedEvent.OldHost.UserId, oldhostId);
            Assert.Equal(hostChangedEvent.NewHost.UserId, newHostId);

            testEnv.Dispose();
        }
        
        [Fact(Timeout = TestsUtil.TIMEOUT_MILLISECONDS)]
        private async void HostShouldBeFirstByAlphanumericSort()
        {
            var testEnv = new SyncTestEnvironment(
                numClients: 4,
                creatorIndex: 0);

            await testEnv.StartAll();
            var expectedHostEnv = testEnv.GetAllUserEnvs().OrderBy(x => x.Self.UserId).First();

            foreach (var env in testEnv.GetAllUserEnvs())
            {
                if (env.Self.UserId == expectedHostEnv.Self.UserId)
                {
                    Assert.True(env.Match.IsSelfHost());
                }
                else
                {
                    Assert.False(env.Match.IsSelfHost());
                }
            }
            
            testEnv.Dispose();
        }
        
        [Fact(Timeout = TestsUtil.TIMEOUT_MILLISECONDS)]
        private async void HostShouldBeSameOnAllClients()
        {
            var testEnv = new SyncTestEnvironment(
                numClients: 4,
                creatorIndex: 0);

            await testEnv.StartAll();
            var expectedHostEnv = testEnv.GetAllUserEnvs().OrderBy(x => x.Self.UserId).First();

            foreach (var env in testEnv.GetAllUserEnvs())
            {
                Assert.Equal(expectedHostEnv.Self.UserId, env.Match.GetHostPresence().UserId);
            }
            
            testEnv.Dispose();
        }
    }
}
